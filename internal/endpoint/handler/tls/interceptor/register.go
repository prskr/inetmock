package interceptor

import (
	"sync"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/inetmock/inetmock/pkg/api"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/pkg/metrics"
	"go.uber.org/zap"
)

var (
	labelNames               = []string{"handler_name"}
	handledRequestCounter    *prometheus.CounterVec
	openConnectionsGauge     *prometheus.GaugeVec
	requestDurationHistogram *prometheus.HistogramVec
)

func AddTLSInterceptor(registry api.HandlerRegistry) (err error) {
	var logger logging.Logger
	if logger, err = logging.CreateLogger(); err != nil {
		panic(err)
	}
	logger = logger.With(
		zap.String("protocol_handler", name),
	)

	if handledRequestCounter, err = metrics.Counter(name, "handled_requests", "", labelNames...); err != nil {
		return
	}
	if openConnectionsGauge, err = metrics.Gauge(name, "open_connections", "", labelNames...); err != nil {
		return
	}
	if requestDurationHistogram, err = metrics.Histogram(name, "request_duration", "", nil, labelNames...); err != nil {

	}

	registry.RegisterHandler(name, func() api.ProtocolHandler {
		return &tlsInterceptor{
			logger:                  logger,
			currentConnectionsCount: new(sync.WaitGroup),
			currentConnections:      make(map[uuid.UUID]*proxyConn),
			connectionsMutex:        &sync.Mutex{},
		}
	})
	return
}
