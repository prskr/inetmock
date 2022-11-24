package proxy

import (
	"github.com/elazarl/goproxy"
	"go.uber.org/zap"

	"inetmock.icb4dc0.de/inetmock/internal/endpoint"
	"inetmock.icb4dc0.de/inetmock/pkg/audit"
	"inetmock.icb4dc0.de/inetmock/pkg/cert"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
)

func AddHTTPProxy(registry endpoint.HandlerRegistry, logger logging.Logger, emitter audit.Emitter, store cert.Store) {
	logger = logger.With(
		zap.String("protocol_handler", name),
	)

	registry.RegisterHandler(name, func() endpoint.ProtocolHandler {
		return &httpProxy{
			logger:    logger,
			emitter:   emitter,
			certStore: store,
			proxy:     goproxy.NewProxyHttpServer(),
		}
	})
}
