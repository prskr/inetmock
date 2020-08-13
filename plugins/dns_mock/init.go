package dns_mock

import (
	"github.com/baez90/inetmock/pkg/api"
	"github.com/baez90/inetmock/pkg/logging"
	"github.com/baez90/inetmock/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
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

func init() {
	var err error
	var logger logging.Logger
	if logger, err = logging.CreateLogger(); err != nil {
		panic(err)
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
		panic(err)
	}

	if unhandledRequestsCounter, err = metrics.Counter(
		name,
		"unhandled_requests_total",
		"",
		handlerNameLblName,
	); err != nil {
		panic(err)
	}

	if requestDurationHistogram, err = metrics.Histogram(
		name,
		"request_duration",
		"",
		nil,
		handlerNameLblName,
	); err != nil {
		panic(err)
	}

	api.Registry().RegisterHandler(name, func() api.ProtocolHandler {
		return &dnsHandler{
			logger: logger,
		}
	})
}
