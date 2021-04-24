package proxy

import (
	"github.com/elazarl/goproxy"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/pkg/metrics"
)

var (
	handlerNameLblName       = "handler_name"
	requestDurationHistogram *prometheus.HistogramVec
)

func AddHTTPProxy(registry endpoint.HandlerRegistry) (err error) {
	var logger logging.Logger
	if logger, err = logging.CreateLogger(); err != nil {
		return
	}
	logger = logger.With(
		zap.String("protocol_handler", name),
	)

	if requestDurationHistogram, err = metrics.Histogram(name, "request_duration", "", nil, handlerNameLblName); err != nil {
		return
	}

	registry.RegisterHandler(name, func() endpoint.ProtocolHandler {
		return &httpProxy{
			logger: logger,
			proxy:  goproxy.NewProxyHttpServer(),
		}
	})

	return
}
