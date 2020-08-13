package http_proxy

import (
	"github.com/baez90/inetmock/pkg/api"
	"github.com/baez90/inetmock/pkg/logging"
	"github.com/baez90/inetmock/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"gopkg.in/elazarl/goproxy.v1"
)

var (
	handlerNameLblName       = "handler_name"
	totalRequestCounter      *prometheus.CounterVec
	totalHttpsRequestCounter *prometheus.CounterVec
	requestDurationHistogram *prometheus.HistogramVec
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

	if totalRequestCounter, err = metrics.Counter(name, "total_requests", "", handlerNameLblName); err != nil {
		panic(err)
	}

	if requestDurationHistogram, err = metrics.Histogram(name, "request_duration", "", nil, handlerNameLblName); err != nil {
		panic(err)
	}

	if totalHttpsRequestCounter, err = metrics.Counter(name, "total_https_requests", "", handlerNameLblName); err != nil {
		panic(err)
	}

	api.Registry().RegisterHandler(name, func() api.ProtocolHandler {
		return &httpProxy{
			logger: logger,
			proxy:  goproxy.NewProxyHttpServer(),
		}
	})
}
