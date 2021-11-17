package pprof

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/pkg/metrics"
)

var (
	totalRequestCounter      *prometheus.CounterVec
	requestDurationHistogram *prometheus.HistogramVec
	initLock                 sync.Locker = new(sync.Mutex)
)

func init() {
	initLock.Lock()
	defer initLock.Unlock()

	var err error
	if totalRequestCounter == nil {
		if totalRequestCounter, err = metrics.Counter(
			name,
			"requests_total",
			"",
		); err != nil {
			panic(err)
		}
	}

	if requestDurationHistogram == nil {
		if requestDurationHistogram, err = metrics.Histogram(
			name,
			"request_duration",
			"",
			nil,
		); err != nil {
			panic(err)
		}
	}
}

func New(logger logging.Logger, emitter audit.Emitter) endpoint.ProtocolHandler {
	return &pprofHandler{
		logger:  logger,
		emitter: emitter,
	}
}

func AddPprof(registry endpoint.HandlerRegistry, logger logging.Logger, emitter audit.Emitter) {
	logger = logger.With(
		zap.String("protocol_handler", name),
	)

	registry.RegisterHandler(name, func() endpoint.ProtocolHandler {
		return New(logger, emitter)
	})
}
