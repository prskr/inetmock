package sink

import (
	"github.com/prometheus/client_golang/prometheus"

	"inetmock.icb4dc0.de/inetmock/pkg/audit"
	"inetmock.icb4dc0.de/inetmock/pkg/metrics"
)

func NewMetricSink() (sink audit.Sink, err error) {
	var totalEventsCounter *prometheus.CounterVec
	if totalEventsCounter, err = metrics.Counter("audit", "events_total", "", "application", "transport"); err != nil {
		return
	}
	sink = &metricSink{
		eventCounter: totalEventsCounter,
	}
	return
}

type metricSink struct {
	eventCounter *prometheus.CounterVec
}

func (metricSink) Name() string {
	return "metrics"
}

func (m metricSink) OnEvent(ev *audit.Event) {
	m.eventCounter.WithLabelValues(ev.Application.String(), ev.Transport.String()).Inc()
}
