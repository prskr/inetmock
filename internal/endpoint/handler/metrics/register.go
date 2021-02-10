package metrics

import (
	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"go.uber.org/zap"
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
