package audit

import "time"

const (
	defaultEventStreamBufferSize  = 100
	defaultSinkBufferSize         = 0
	defaultSinkConsumptionTimeout = 50 * time.Millisecond
)

var (
	generatorIdx   int64 = 1
	WithBufferSize       = func(bufferSize int) EventStreamOption {
		return func(cfg *eventStreamCfg) {
			cfg.bufferSize = bufferSize
		}
	}
	WithGeneratorIndex = func(generatorIndex int64) EventStreamOption {
		return func(cfg *eventStreamCfg) {
			cfg.generatorIndex = generatorIndex
		}
	}
	WithSinkBufferSize = func(bufferSize int) EventStreamOption {
		return func(cfg *eventStreamCfg) {
			cfg.sinkBuffersize = bufferSize
		}
	}
	WithSinkConsumptionTimeout = func(timeout time.Duration) EventStreamOption {
		return func(cfg *eventStreamCfg) {
			cfg.sinkConsumptionTimeout = timeout
		}
	}
)

type EventStreamOption func(cfg *eventStreamCfg)

type eventStreamCfg struct {
	bufferSize             int
	sinkBuffersize         int
	generatorIndex         int64
	sinkConsumptionTimeout time.Duration
}

func newEventStreamCfg() eventStreamCfg {
	cfg := eventStreamCfg{
		generatorIndex:         generatorIdx,
		sinkBuffersize:         defaultSinkBufferSize,
		bufferSize:             defaultEventStreamBufferSize,
		sinkConsumptionTimeout: defaultSinkConsumptionTimeout,
	}
	generatorIdx++

	return cfg
}
