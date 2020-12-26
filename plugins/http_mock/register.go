package http_mock

import (
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/inetmock/inetmock/pkg/api"
	"gitlab.com/inetmock/inetmock/pkg/metrics"
)

var (
	totalRequestCounter      *prometheus.CounterVec
	requestDurationHistogram *prometheus.HistogramVec
)

func AddHTTPMock(registry api.HandlerRegistry) (err error) {
	if totalRequestCounter == nil {
		if totalRequestCounter, err = metrics.Counter(
			name,
			"total_requests",
			"",
			handlerNameLblName,
			ruleMatchedLblName,
		); err != nil {
			return
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
			return
		}
	}

	registry.RegisterHandler(name, func() api.ProtocolHandler {
		return &httpHandler{}
	})

	return
}
