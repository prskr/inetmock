package audit

import (
	"context"
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"go.uber.org/zap"
)

func init() {
	snowflake.Epoch = time.Unix(0, 0).Unix()
}

type eventStream struct {
	logger                 logging.Logger
	buffer                 chan *Event
	sinks                  map[string]*registeredSink
	lock                   sync.Locker
	idGenerator            *snowflake.Node
	sinkBufferSize         int
	sinkConsumptionTimeout time.Duration
}

type registeredSink struct {
	downstream chan Event
	lock       sync.Locker
}

func NewEventStream(logger logging.Logger, options ...EventStreamOption) (es EventStream, err error) {
	cfg := newEventStreamCfg()

	for _, opt := range options {
		opt(&cfg)
	}

	var node *snowflake.Node
	node, err = snowflake.NewNode(cfg.generatorIndex)
	if err != nil {
		return
	}

	generatorIdx++
	underlying := &eventStream{
		logger:                 logger,
		sinks:                  make(map[string]*registeredSink),
		buffer:                 make(chan *Event, cfg.bufferSize),
		sinkBufferSize:         cfg.sinkBuffersize,
		sinkConsumptionTimeout: cfg.sinkConsumptionTimeout,
		idGenerator:            node,
		lock:                   &sync.Mutex{},
	}

	// start distribute workers
	for i := 0; i < cfg.distributeParallelization; i++ {
		go underlying.distribute()
	}

	es = underlying

	return
}

func (e *eventStream) Emit(ev Event) {
	ev.ApplyDefaults(e.idGenerator.Generate().Int64())
	select {
	case e.buffer <- &ev:
		e.logger.Debug("pushed event to distribute loop")
	default:
		e.logger.Warn("buffer is full")
	}
}

func (e *eventStream) RemoveSink(name string) (exists bool) {
	e.lock.Lock()
	defer e.lock.Unlock()

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
	name := s.Name()

	if _, exists := e.sinks[name]; exists {
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
	for name := range e.sinks {
		sinks = append(sinks, name)
	}
	return
}

func (e *eventStream) Close() error {
	close(e.buffer)
	for _, rs := range e.sinks {
		close(rs.downstream)
	}
	return nil
}

func (e *eventStream) distribute() {
	for ev := range e.buffer {
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
	}
}
