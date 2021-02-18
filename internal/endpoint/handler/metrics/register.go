package metrics

import (
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

func AddMetricsExporter(registry endpoint.HandlerRegistry) (err error) {
	var logger logging.Logger
	if logger, err = logging.CreateLogger(); err != nil {
		return
	}
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
