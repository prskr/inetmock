package dns_mock

import (
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/inetmock/inetmock/pkg/api"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/pkg/metrics"
	"go.uber.org/zap"
)

const (
	name = "dns_mock"
)

var (
	handlerNameLblName          = "handler_name"
	totalHandledRequestsCounter *prometheus.CounterVec
	unhandledRequestsCounter    *prometheus.CounterVec
	requestDurationHistogram    *prometheus.HistogramVec
)

func AddDNSMock(registry api.HandlerRegistry) (err error) {
	var logger logging.Logger
	if logger, err = logging.CreateLogger(); err != nil {
		return
	}
	logger = logger.With(
		zap.String("protocol_handler", name),
	)

	if totalHandledRequestsCounter, err = metrics.Counter(
		name,
		"handled_requests_total",
		"",
		handlerNameLblName,
	); err != nil {
		return
	}

	if unhandledRequestsCounter, err = metrics.Counter(
		name,
		"unhandled_requests_total",
		"",
		handlerNameLblName,
	); err != nil {
		return
	}

	if requestDurationHistogram, err = metrics.Histogram(
		name,
		"request_duration",
		"",
		nil,
		handlerNameLblName,
	); err != nil {
		return
	}

	registry.RegisterHandler(name, func() api.ProtocolHandler {
		return &dnsHandler{
			logger: logger,
		}
	})

	return
}
