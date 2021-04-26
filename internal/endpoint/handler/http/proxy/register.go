package proxy

import (
	"github.com/elazarl/goproxy"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/cert"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/pkg/metrics"
)

var (
	handlerNameLblName       = "handler_name"
	requestDurationHistogram *prometheus.HistogramVec
)

func AddHTTPProxy(registry endpoint.HandlerRegistry, logger logging.Logger, emitter audit.Emitter, store cert.Store) (err error) {
	logger = logger.With(
		zap.String("protocol_handler", name),
	)

	if requestDurationHistogram, err = metrics.Histogram(name, "request_duration", "", nil, handlerNameLblName); err != nil {
		return
	}

	registry.RegisterHandler(name, func() endpoint.ProtocolHandler {
		return &httpProxy{
			logger:    logger,
			emitter:   emitter,
			certStore: store,
			proxy:     goproxy.NewProxyHttpServer(),
		}
	})

	return
}
