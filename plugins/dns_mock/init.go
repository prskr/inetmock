package dns_mock

import (
	"github.com/baez90/inetmock/pkg/api"
	"github.com/baez90/inetmock/pkg/logging"
	"go.uber.org/zap"
)

const (
	name = "dns_mock"
)

func init() {
	logger, _ := logging.CreateLogger()
	logger = logger.With(
		zap.String("ProtocolHandler", name),
	)

	api.Registry().RegisterHandler(name, func() api.ProtocolHandler {
		return &dnsHandler{
			logger: logger,
		}
	})
}
