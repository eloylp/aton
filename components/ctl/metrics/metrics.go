package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "aton"
	subsystem = "ctl"
)

func currentDetectors() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "detector_up",
		Help:      "Detectors that are up",
	}, []string{"uuid"})
}
