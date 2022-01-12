package dhcp

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/internal/state"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/pkg/metrics"
)

var (
	requestDurationHistogram *prometheus.HistogramVec
	initLock                 sync.Mutex
)

func init() {
	initLock.Lock()
	defer initLock.Unlock()

	var err error
	if requestDurationHistogram == nil {
		if requestDurationHistogram, err = metrics.Histogram(
			name,
			"request_duration",
			"",
			nil,
			handlerNameLblName,
		); err != nil {
			panic(err)
		}
	}
}

func New(logger logging.Logger, emitter audit.Emitter, stateStore state.KVStore) endpoint.ProtocolHandler {
	return &dhcpHandler{
		logger:     logger,
		emitter:    emitter,
		stateStore: stateStore,
	}
}

func AddDHCPMock(registry endpoint.HandlerRegistry, logger logging.Logger, emitter audit.Emitter, stateStore state.KVStore) {
	registry.RegisterHandler(name, func() endpoint.ProtocolHandler {
		return New(logger, emitter, stateStore)
	})
}
