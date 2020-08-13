package metrics_exporter

import (
	"context"
	"errors"
	"github.com/baez90/inetmock/pkg/config"
	"github.com/baez90/inetmock/pkg/logging"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
)

const (
	name = "metrics_exporter"
)

type metricsExporter struct {
	logger logging.Logger
	server *http.Server
}

func (m *metricsExporter) Start(config config.HandlerConfig) (err error) {
	exporterOptions := metricsExporterOptions{}
	if err = config.Options.Unmarshal(&exporterOptions); err != nil {
		return
	}

	m.logger = m.logger.With(
		zap.String("handler_name", config.HandlerName),
		zap.String("address", config.ListenAddr()),
	)

	mux := http.NewServeMux()
	mux.Handle(exporterOptions.Route, promhttp.Handler())
	m.server = &http.Server{
		Addr:    config.ListenAddr(),
		Handler: mux,
	}

	go func() {
		if err := m.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			m.logger.Error(
				"Error occurred while serving metrics",
				zap.Error(err),
			)
		}
	}()
	return
}

func (m *metricsExporter) Shutdown(ctx context.Context) error {
	return m.server.Shutdown(ctx)
}
