package mock

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
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

func AddHTTPMock(registry endpoint.HandlerRegistry) (err error) {
	if err := InitMetrics(); err != nil {
		return err
	}

	registry.RegisterHandler(name, func() endpoint.ProtocolHandler {
		return &httpHandler{}
	})

	return
}
