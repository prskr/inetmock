package protocols

import (
	"github.com/prometheus/client_golang/prometheus"

	"inetmock.icb4dc0.de/inetmock/pkg/metrics"
)

var RequestDurationHistogram *prometheus.HistogramVec

func init() {
	var err error
	if RequestDurationHistogram == nil {
		if RequestDurationHistogram, err = metrics.Histogram(
			"protocols",
			"request_duration",
			"",
			nil,
			"protocol",
			"handler_name",
		); err != nil {
			panic(err)
		}
	}
}
