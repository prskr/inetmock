package netflow

import (
	"github.com/prometheus/client_golang/prometheus"

	"inetmock.icb4dc0.de/inetmock/pkg/metrics"
)

var connTrackGauge *prometheus.GaugeVec

func init() {
	var err error
	connTrackGauge, err = metrics.Gauge(
		"netflow",
		"conn_track_entries",
		`Describe how many entries are currently stored in conn_track map`,
		"interface",
	)

	if err != nil {
		panic(err)
	}
}
