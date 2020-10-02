package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	metricNamespace = "inetmock"
)

func Gauge(subsystem, name, help string, labelNames ...string) (*prometheus.GaugeVec, error) {
	vec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: metricNamespace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
		}, labelNames)
	return vec, prometheus.Register(vec)
}

func Histogram(subsystem, name, help string, buckets []float64, labelNames ...string) (*prometheus.HistogramVec, error) {
	vec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: metricNamespace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
			Buckets:   buckets,
		},
		labelNames,
	)
	return vec, prometheus.Register(vec)
}

func Counter(subsystem, name, help string, labelNames ...string) (*prometheus.CounterVec, error) {
	vec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricNamespace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
		},
		labelNames,
	)
	return vec, prometheus.Register(vec)
}
