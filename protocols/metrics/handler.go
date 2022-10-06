package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/pkg/health"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

const (
	name        = "metrics_exporter"
	healthRoute = "/health"
)

type metricsExporter struct {
	logger  logging.Logger
	checker health.Checker
	server  *http.Server
}

func (m *metricsExporter) Start(_ context.Context, startupSpec *endpoint.StartupSpec) error {
	var exporterOptions metricsExporterOptions
	if err := startupSpec.UnmarshalOptions(&exporterOptions); err != nil {
		return err
	}

	m.logger = m.logger.With(
		zap.String("handler_name", startupSpec.Name),
		zap.String("address", startupSpec.Addr.String()),
	)

	mux := http.NewServeMux()
	mux.Handle(exporterOptions.Route, promhttp.Handler())
	mux.Handle(healthRoute, health.NewHealthHandler(m.checker))

	m.server = &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 50 * time.Millisecond,
	}

	go func() {
		if err := endpoint.IgnoreShutdownError(m.server.Serve(startupSpec.Listener)); err != nil {
			m.logger.Error(
				"Error occurred while serving metrics",
				zap.Error(err),
			)
		}
	}()

	return nil
}
