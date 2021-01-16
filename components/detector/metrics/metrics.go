package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "aton"
const subsystem = "detector"

// Capturer (input video stream) metrics

func capturerReceivedFramesTotal() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "capturer_received_frames_total",
		Help:      "The total frames obtained from capturers",
	}, []string{"uuid", "capturer_uuid"})
}

func capturerFailedFramesTotal() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "capturer_failed_frames_total",
		Help:      "The failed frames returned by capturers",
	}, []string{"uuid", "capturer_uuid"})
}

// Responses of detectors metrics

func processedFramesTotal() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "processed_frames_total",
		Help:      "The total frame processed detectors",
	}, []string{"uuid"})
}

func failedFramesTotal() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "error_frames_total",
		Help:      "The failed, not processed, frames received from detectors",
	}, []string{"uuid"})
}

func unrecognizedFramesTotal() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "unrecognized_frames_total",
		Help:      "The total unrecognized frames processed by detectors",
	}, []string{"uuid"})
}

// Status gauges metrics. Provides information about running systems

func currentCapturers() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "capturer_up",
		Help:      "Capturers that are up",
	}, []string{"uuid"})
}
