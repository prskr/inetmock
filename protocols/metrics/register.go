package metrics

import (
	"go.uber.org/zap"

	"inetmock.icb4dc0.de/inetmock/internal/endpoint"
	"inetmock.icb4dc0.de/inetmock/pkg/health"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
)

func AddMetricsExporter(registry endpoint.HandlerRegistry, logger logging.Logger, checker health.Checker) {
	logger = logger.With(
		zap.String("protocol_handler", name),
	)

	registry.RegisterHandler(name, func() endpoint.ProtocolHandler {
		return &metricsExporter{
			logger:  logger,
			checker: checker,
		}
	})
}
