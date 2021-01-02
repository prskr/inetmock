package audit

import (
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
	sinks                  map[string]chan Event
	lock                   sync.Locker
	idGenerator            *snowflake.Node
	sinkBufferSize         int
	sinkConsumptionTimeout time.Duration
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
		sinks:                  make(map[string]chan Event),
		buffer:                 make(chan *Event, cfg.bufferSize),
		sinkBufferSize:         cfg.sinkBuffersize,
		sinkConsumptionTimeout: cfg.sinkConsumptionTimeout,
		idGenerator:            node,
		lock:                   &sync.Mutex{},
	}

	go underlying.distribute()

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

func (e *eventStream) RegisterSink(s Sink) error {
	name := s.Name()

	if _, exists := e.sinks[name]; exists {
		return ErrSinkAlreadyRegistered
	}

	downstream := make(chan Event, e.sinkBufferSize)
	s.OnSubscribe(downstream)
	e.sinks[name] = downstream
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
	for _, downstream := range e.sinks {
		close(downstream)
	}
	return nil
}

func (e *eventStream) distribute() {
	for ev := range e.buffer {
		for name, s := range e.sinks {
			e.logger.Debug("notify sink", zap.String("sink", name))
			select {
			case s <- *ev:
				e.logger.Debug("pushed event to sink channel")
			case <-time.After(e.sinkConsumptionTimeout):
				e.logger.Warn("sink consummation timed out")
			}
		}
	}
}
