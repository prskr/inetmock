package metrics

import (
	"context"
	"errors"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

const (
	name = "metrics_exporter"
)

type metricsExporter struct {
	logger logging.Logger
	server *http.Server
}

func (m *metricsExporter) Start(ctx context.Context, lifecycle endpoint.Lifecycle) error {
	var exporterOptions metricsExporterOptions
	if err := lifecycle.UnmarshalOptions(&exporterOptions); err != nil {
		return err
	}

	m.logger = m.logger.With(
		zap.String("handler_name", lifecycle.Name()),
		zap.String("address", lifecycle.Uplink().Addr().String()),
	)

	mux := http.NewServeMux()
	mux.Handle(exporterOptions.Route, promhttp.Handler())
	m.server = &http.Server{
		Handler: mux,
	}

	go func() {
		if err := m.server.Serve(lifecycle.Uplink().Listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			m.logger.Error(
				"Error occurred while serving metrics",
				zap.Error(err),
			)
		}
	}()

	go func() {
		<-ctx.Done()
		if err := m.server.Close(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			m.logger.Error("failed to stop metrics server", zap.Error(err))
		}
	}()
	return nil
}
