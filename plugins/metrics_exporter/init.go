package metrics_exporter

import (
	"github.com/baez90/inetmock/pkg/api"
	"github.com/baez90/inetmock/pkg/logging"
	"go.uber.org/zap"
)

func init() {
	logger, _ := logging.CreateLogger()
	logger = logger.With(
		zap.String("protocol_handler", name),
	)

	api.Registry().RegisterHandler(name, func() api.ProtocolHandler {
		return &metricsExporter{
			logger: logger,
		}
	})
}
