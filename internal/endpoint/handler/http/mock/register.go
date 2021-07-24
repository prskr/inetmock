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

func init() {
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
			panic(err)
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
			panic(err)
		}
	}
}

func New(logger logging.Logger, emitter audit.Emitter, fakeFileFS fs.FS) endpoint.ProtocolHandler {
	return &httpHandler{
		logger:     logger,
		fakeFileFS: fakeFileFS,
		emitter:    emitter,
	}
}

func AddHTTPMock(registry endpoint.HandlerRegistry, logger logging.Logger, emitter audit.Emitter, fakeFileFS fs.FS) {
	registry.RegisterHandler(name, func() endpoint.ProtocolHandler {
		return New(logger, emitter, fakeFileFS)
	})
}
