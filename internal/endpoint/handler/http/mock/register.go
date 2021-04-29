package mock

import (
	"io/fs"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

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

func InitMetrics() error {
	initLock.Lock()
	defer initLock.Unlock()

	var err error
	if totalRequestCounter == nil {
		if totalRequestCounter, err = metrics.Counter(
			name,
			"total_requests",
			"",
			handlerNameLblName,
			ruleMatchedLblName,
		); err != nil {
			return err
		}
	}

	if requestDurationHistogram == nil {
		if requestDurationHistogram, err = metrics.Histogram(
			name,
			"request_duration",
			"",
			nil,
			handlerNameLblName,
		); err != nil {
			return err
		}
	}
	return nil
}

func New(logger logging.Logger, emitter audit.Emitter, fakeFileFS fs.FS) endpoint.ProtocolHandler {
	return &httpHandler{
		logger:     logger,
		fakeFileFS: fakeFileFS,
		emitter:    emitter,
	}
}

func AddHTTPMock(registry endpoint.HandlerRegistry, logger logging.Logger, emitter audit.Emitter, fakeFileFS fs.FS) (err error) {
	if err := InitMetrics(); err != nil {
		return err
	}

	registry.RegisterHandler(name, func() endpoint.ProtocolHandler {
		return New(logger, emitter, fakeFileFS)
	})

	return
}
