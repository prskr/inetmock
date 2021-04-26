package metrics

import (
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

func AddMetricsExporter(registry endpoint.HandlerRegistry, logger logging.Logger) (err error) {
	logger = logger.With(
		zap.String("protocol_handler", name),
	)

	registry.RegisterHandler(name, func() endpoint.ProtocolHandler {
		return &metricsExporter{
			logger: logger,
		}
	})
	return
}
