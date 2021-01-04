package audit

import (
	"runtime"
	"time"
)

const (
	defaultEventStreamBufferSize  = 100
	defaultSinkBufferSize         = 0
	defaultSinkConsumptionTimeout = 50 * time.Millisecond
)

var (
	generatorIdx                     int64 = 1
	defaultDistributeParallelization       = runtime.NumCPU() / 2
	WithBufferSize                         = func(bufferSize int) EventStreamOption {
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
	WithDistributeParallelization = func(parallelization int) EventStreamOption {
		return func(cfg *eventStreamCfg) {
			if parallelization <= 0 || parallelization > runtime.NumCPU() {
				return
			}
			cfg.distributeParallelization = parallelization
		}
	}
)

type EventStreamOption func(cfg *eventStreamCfg)

type eventStreamCfg struct {
	bufferSize                int
	sinkBuffersize            int
	generatorIndex            int64
	distributeParallelization int
	sinkConsumptionTimeout    time.Duration
}

func newEventStreamCfg() eventStreamCfg {
	cfg := eventStreamCfg{
		generatorIndex:            generatorIdx,
		sinkBuffersize:            defaultSinkBufferSize,
		bufferSize:                defaultEventStreamBufferSize,
		sinkConsumptionTimeout:    defaultSinkConsumptionTimeout,
		distributeParallelization: defaultDistributeParallelization,
	}
	generatorIdx++

	return cfg
}
