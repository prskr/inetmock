package metrics_exporter

import (
	"gitlab.com/inetmock/inetmock/pkg/api"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"go.uber.org/zap"
)

func AddMetricsExporter(registry api.HandlerRegistry) (err error) {
	var logger logging.Logger
	if logger, err = logging.CreateLogger(); err != nil {
		return
	}
	logger = logger.With(
		zap.String("protocol_handler", name),
	)

	registry.RegisterHandler(name, func() api.ProtocolHandler {
		return &metricsExporter{
			logger: logger,
		}
	})
	return
}
