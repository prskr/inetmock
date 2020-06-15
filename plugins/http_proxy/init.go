package http_proxy

import (
	"github.com/baez90/inetmock/pkg/api"
	"github.com/baez90/inetmock/pkg/logging"
	"go.uber.org/zap"
	"gopkg.in/elazarl/goproxy.v1"
)

func init() {
	logger, _ := logging.CreateLogger()
	logger = logger.With(
		zap.String("ProtocolHandler", name),
	)

	api.Registry().RegisterHandler(name, func() api.ProtocolHandler {
		return &httpProxy{
			logger: logger,
			proxy:  goproxy.NewProxyHttpServer(),
		}
	})
}
