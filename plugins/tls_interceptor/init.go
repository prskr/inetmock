package tls_interceptor

import (
	"github.com/baez90/inetmock/pkg/api"
	"github.com/baez90/inetmock/pkg/logging"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"sync"
)

func init() {
	logger, _ := logging.CreateLogger()
	logger = logger.With(
		zap.String("ProtocolHandler", name),
	)

	api.Registry().RegisterHandler(name, func() api.ProtocolHandler {
		return &tlsInterceptor{
			logger:                  logger,
			currentConnectionsCount: &sync.WaitGroup{},
			currentConnections:      make(map[uuid.UUID]*proxyConn),
			connectionsMutex:        &sync.Mutex{},
		}
	})
}
