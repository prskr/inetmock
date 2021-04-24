package mock

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/pkg/metrics"
)

const (
	name = "dns_mock"
)

var (
	handlerNameLblName          = "handler_name"
	totalHandledRequestsCounter *prometheus.CounterVec
	unhandledRequestsCounter    *prometheus.CounterVec
	requestDurationHistogram    *prometheus.HistogramVec
	initLock                    sync.Locker = new(sync.Mutex)
)

func AddDNSMock(registry endpoint.HandlerRegistry) error {
	initLock.Lock()
	defer initLock.Unlock()

	var err error
	if totalHandledRequestsCounter == nil {
		if totalHandledRequestsCounter, err = metrics.Counter(
			name,
			"handled_requests_total",
			"",
			handlerNameLblName,
		); err != nil {
			return err
		}
	}

	if unhandledRequestsCounter == nil {
		if unhandledRequestsCounter, err = metrics.Counter(
			name,
			"unhandled_requests_total",
			"",
			handlerNameLblName,
		); err != nil {
			return err
		}
	}

	if requestDurationHistogram == nil {
		if requestDurationHistogram, err = metrics.Histogram(
			name,
			"request_duration",
			"",
			nil,
			handlerNameLblName,
		); err != nil {
			return err
		}
	}

	registry.RegisterHandler(name, func() endpoint.ProtocolHandler {
		return &dnsHandler{}
	})

	return nil
}
