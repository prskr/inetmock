package tls_interceptor

import (
	"github.com/baez90/inetmock/pkg/api"
	"github.com/baez90/inetmock/pkg/logging"
	"github.com/baez90/inetmock/pkg/metrics"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"sync"
)

var (
	labelNames               = []string{"handler_name"}
	handledRequestCounter    *prometheus.CounterVec
	openConnectionsGauge     *prometheus.GaugeVec
	requestDurationHistogram *prometheus.HistogramVec
)

func init() {
	var err error
	var logger logging.Logger
	if logger, err = logging.CreateLogger(); err != nil {
		panic(err)
	}
	logger = logger.With(
		zap.String("protocol_handler", name),
	)

	if handledRequestCounter, err = metrics.Counter(name, "handled_requests", "", labelNames...); err != nil {
		panic(err)
	}
	if openConnectionsGauge, err = metrics.Gauge(name, "open_connections", "", labelNames...); err != nil {
		panic(err)
	}
	if requestDurationHistogram, err = metrics.Histogram(name, "request_duration", "", nil, labelNames...); err != nil {
		panic(err)
	}

	api.Registry().RegisterHandler(name, func() api.ProtocolHandler {
		return &tlsInterceptor{
			logger:                  logger,
			currentConnectionsCount: &sync.WaitGroup{},
			currentConnections:      make(map[uuid.UUID]*proxyConn),
			connectionsMutex:        &sync.Mutex{},
		}
	})
}
