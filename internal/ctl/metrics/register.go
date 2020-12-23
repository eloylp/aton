package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Register struct {
	capturerReceivedFramesTotal *prometheus.CounterVec
	capturerFailedFramesTotal   *prometheus.CounterVec
	processedFramesTotal        *prometheus.CounterVec
	failedFramesTotal           *prometheus.CounterVec
	currentCapturers            *prometheus.GaugeVec
	currentDetectors            *prometheus.GaugeVec
}

func NewRegister() *Register {
	return &Register{
		capturerReceivedFramesTotal: capturerReceivedFramesTotal,
		capturerFailedFramesTotal:   capturerFailedFramesTotal,
		processedFramesTotal:        processedFramesTotal,
		failedFramesTotal:           failedFramesTotal,
		currentCapturers:            currentCapturers,
		currentDetectors:            currentDetectors,
	}
}

func (r *Register) IncCapturerReceivedFramesTotal(labelValues ...string) {
	r.capturerReceivedFramesTotal.WithLabelValues(labelValues...).Inc()
}

func (r *Register) IncCapturerFailedFramesTotal(labelValues ...string) {
	r.capturerFailedFramesTotal.WithLabelValues(labelValues...).Inc()
}

func (r *Register) IncProcessedFramesTotal(labelValues ...string) {
	r.processedFramesTotal.WithLabelValues(labelValues...).Inc()
}

func (r *Register) IncFailedFramesTotal(labelValues ...string) {
	r.failedFramesTotal.WithLabelValues(labelValues...).Inc()
}

func (r *Register) CapturerUP(labelValues ...string) {
	r.currentCapturers.WithLabelValues(labelValues...).Inc()
}

func (r *Register) CapturerDown(labelValues ...string) {
	r.currentCapturers.WithLabelValues(labelValues...).Dec()
}

func (r *Register) DetectorUP(labelValues ...string) {
	r.currentDetectors.WithLabelValues(labelValues...).Inc()
}

func (r *Register) DetectorDown(labelValues ...string) {
	r.currentDetectors.WithLabelValues(labelValues...).Dec()
}
