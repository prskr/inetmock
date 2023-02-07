package audit

import (
	"context"
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"

	"inetmock.icb4dc0.de/inetmock/pkg/logging"
)

const (
	emitTimeout = 10 * time.Millisecond
)

var _ EventStream = (*eventStream)(nil)

func init() {
	snowflake.Epoch = time.Unix(0, 0).Unix()
}

type eventStream struct {
	logger                 logging.Logger
	buffer                 chan *Event
	sinks                  map[string]*lockableSink
	rwlock                 sync.RWMutex
	idGenerator            *snowflake.Node
	sinkBufferSize         int
	sinkConsumptionTimeout time.Duration
}

type lockableSink struct {
	Sink
	lock sync.Mutex
}

func (s *lockableSink) OnEvent(ev *Event) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Sink.OnEvent(ev)
}

func NewEventStream(logger logging.Logger, options ...EventStreamOption) (EventStream, error) {
	cfg := newEventStreamCfg()

	for _, opt := range options {
		opt(&cfg)
	}

	var err error
	var node *snowflake.Node
	if node, err = snowflake.NewNode(cfg.generatorIndex); err != nil {
		return nil, err
	}

	atomic.AddInt64(&generatorIdx, 1)
	underlying := &eventStream{
		logger:                 logger,
		sinks:                  make(map[string]*lockableSink),
		buffer:                 make(chan *Event, cfg.bufferSize),
		sinkBufferSize:         cfg.sinkBuffersize,
		sinkConsumptionTimeout: cfg.sinkConsumptionTimeout,
		idGenerator:            node,
	}

	// start distribute workers
	for i := 0; i < cfg.distributeParallelization; i++ {
		go underlying.distribute()
	}

	return underlying, err
}

func (e *eventStream) Emit(ev *Event) {
	ev.ApplyDefaults(e.idGenerator.Generate().Int64())
	select {
	case e.buffer <- ev:
		e.logger.Debug("pushed event to distribute loop")
	case <-time.After(emitTimeout):
		e.logger.Warn("buffer is full")
	}
}

func (e *eventStream) Builder() EventBuilder {
	return BuilderForEmitter(e)
}

func (e *eventStream) RemoveSink(name string) (exists bool) {
	e.rwlock.Lock()
	defer e.rwlock.Unlock()

	var sink *lockableSink
	sink, exists = e.sinks[name]
	if !exists {
		return
	}
	sink.lock.Lock()
	defer sink.lock.Unlock()
	delete(e.sinks, name)

	return
}

func (e *eventStream) RegisterSink(ctx context.Context, s Sink) error {
	e.rwlock.Lock()
	defer e.rwlock.Unlock()

	name := s.Name()
	if _, present := e.sinks[name]; present {
		return ErrSinkAlreadyRegistered
	}

	rs := &lockableSink{
		Sink: s,
	}

	go func() {
		<-ctx.Done()
		e.RemoveSink(name)
	}()

	e.sinks[name] = rs
	return nil
}

func (e *eventStream) Sinks() (sinks []string) {
	e.rwlock.RLock()
	defer e.rwlock.RUnlock()

	return maps.Keys(e.sinks)
}

func (e *eventStream) Close() error {
	e.rwlock.Lock()
	defer e.rwlock.Unlock()

	close(e.buffer)
	var err error
	for _, rs := range e.sinks {
		if closer, ok := rs.Sink.(io.Closer); ok {
			err = errors.Join(err, closer.Close())
		}
	}
	return nil
}

func (e *eventStream) distribute() {
	var wg sync.WaitGroup
	for ev := range e.buffer {
		e.rwlock.RLock()
		wg.Add(len(e.sinks))
		for name, rs := range e.sinks {
			go func(name string, rs *lockableSink, wg *sync.WaitGroup) {
				e.logger.Debug("notify sink", zap.String("sink", name))
				rs.OnEvent(ev)
				wg.Done()
			}(name, rs, &wg)
		}

		wg.Wait()
		ev.Dispose()
		e.rwlock.RUnlock()
	}
}
