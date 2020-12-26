package http_proxy

import (
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/inetmock/inetmock/pkg/api"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/pkg/metrics"
	"go.uber.org/zap"
	"gopkg.in/elazarl/goproxy.v1"
)

var (
	handlerNameLblName       = "handler_name"
	totalRequestCounter      *prometheus.CounterVec
	totalHttpsRequestCounter *prometheus.CounterVec
	requestDurationHistogram *prometheus.HistogramVec
)

func AddHTTPProxy(registry api.HandlerRegistry) (err error) {
	var logger logging.Logger
	if logger, err = logging.CreateLogger(); err != nil {
		return
	}
	logger = logger.With(
		zap.String("protocol_handler", name),
	)

	if totalRequestCounter, err = metrics.Counter(name, "total_requests", "", handlerNameLblName); err != nil {
		return
	}

	if requestDurationHistogram, err = metrics.Histogram(name, "request_duration", "", nil, handlerNameLblName); err != nil {
		return
	}

	if totalHttpsRequestCounter, err = metrics.Counter(name, "total_https_requests", "", handlerNameLblName); err != nil {
		return
	}

	registry.RegisterHandler(name, func() api.ProtocolHandler {
		return &httpProxy{
			logger: logger,
			proxy:  goproxy.NewProxyHttpServer(),
		}
	})

	return
}
