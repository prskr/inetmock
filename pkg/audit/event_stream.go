package audit

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/pkg/logging"
)

const (
	emitTimeout = 10 * time.Millisecond
)

func init() {
	snowflake.Epoch = time.Unix(0, 0).Unix()
}

type eventStream struct {
	logger                 logging.Logger
	buffer                 chan *Event
	sinks                  map[string]*registeredSink
	readLock               sync.Locker
	writeLock              sync.Locker
	idGenerator            *snowflake.Node
	sinkBufferSize         int
	sinkConsumptionTimeout time.Duration
}

type registeredSink struct {
	downstream chan Event
	lock       sync.Locker
}

func MustNewEventStream(logger logging.Logger, options ...EventStreamOption) EventStream {
	if stream, err := NewEventStream(logger, options...); err != nil {
		panic(err)
	} else {
		return stream
	}
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

	rwMutex := new(sync.RWMutex)
	atomic.AddInt64(&generatorIdx, 1)
	underlying := &eventStream{
		logger:                 logger,
		sinks:                  make(map[string]*registeredSink),
		buffer:                 make(chan *Event, cfg.bufferSize),
		sinkBufferSize:         cfg.sinkBuffersize,
		sinkConsumptionTimeout: cfg.sinkConsumptionTimeout,
		idGenerator:            node,
		writeLock:              rwMutex,
		readLock:               rwMutex.RLocker(),
	}

	// start distribute workers
	for i := 0; i < cfg.distributeParallelization; i++ {
		go underlying.distribute()
	}

	return underlying, err
}

func (e *eventStream) Emit(ev Event) {
	ev.ApplyDefaults(e.idGenerator.Generate().Int64())
	select {
	case e.buffer <- &ev:
		e.logger.Debug("pushed event to distribute loop")
	case <-time.After(emitTimeout):
		e.logger.Warn("buffer is full")
	}
}

func (e *eventStream) RemoveSink(name string) (exists bool) {
	e.writeLock.Lock()
	defer e.writeLock.Unlock()

	var sink *registeredSink
	sink, exists = e.sinks[name]
	if !exists {
		return
	}
	sink.lock.Lock()
	defer sink.lock.Unlock()
	delete(e.sinks, name)
	close(sink.downstream)

	return
}

func (e *eventStream) RegisterSink(ctx context.Context, s Sink) error {
	e.writeLock.Lock()
	defer e.writeLock.Unlock()

	name := s.Name()
	if _, present := e.sinks[name]; present {
		return ErrSinkAlreadyRegistered
	}

	rs := &registeredSink{
		downstream: make(chan Event, e.sinkBufferSize),
		lock:       new(sync.Mutex),
	}

	s.OnSubscribe(rs.downstream)

	go func() {
		<-ctx.Done()
		e.RemoveSink(name)
	}()

	e.sinks[name] = rs
	return nil
}

func (e eventStream) Sinks() (sinks []string) {
	e.readLock.Lock()
	defer e.readLock.Unlock()

	for name := range e.sinks {
		sinks = append(sinks, name)
	}
	return
}

func (e *eventStream) Close() error {
	e.writeLock.Lock()
	defer e.writeLock.Unlock()

	close(e.buffer)
	for _, rs := range e.sinks {
		close(rs.downstream)
	}
	return nil
}

func (e *eventStream) distribute() {
	for ev := range e.buffer {
		e.readLock.Lock()
		for name, rs := range e.sinks {
			rs.lock.Lock()
			e.logger.Debug("notify sink", zap.String("sink", name))
			select {
			case rs.downstream <- *ev:
				e.logger.Debug("pushed event to sink channel")
			case <-time.After(e.sinkConsumptionTimeout):
				e.logger.Warn("sink consummation timed out")
			}
			rs.lock.Unlock()
		}
		e.readLock.Unlock()
	}
}
