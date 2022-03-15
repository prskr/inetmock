package sink

import (
	"github.com/prometheus/client_golang/prometheus"

	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/metrics"
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

func (m metricSink) OnSubscribe(evs <-chan *audit.Event) {
	go func(evs <-chan *audit.Event) {
		for ev := range evs {
			m.eventCounter.WithLabelValues(ev.Application.String(), ev.Transport.String()).Inc()
		}
	}(evs)
}
