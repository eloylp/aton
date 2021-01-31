package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "aton"
	subsystem = "detector"
)

// Capturer (input video stream) metrics

func capturerReceivedFramesTotal() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "capturer_received_frames_total",
		Help:      "The total frames obtained from capturers",
	}, []string{"uuid", "capturer_uuid", "capturer_url"})
}

func capturerFailedFramesTotal() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "capturer_failed_frames_total",
		Help:      "The failed frames returned by capturers",
	}, []string{"uuid", "capturer_uuid", "capturer_url"})
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

func entitiesTotal() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "entities_total",
		Help:      "The total entities processed by detectors",
	}, []string{"uuid"})
}

func unrecognizedEntitiesTotal() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "unrecognized_entities_total",
		Help:      "The total unrecognized entities processed by detectors",
	}, []string{"uuid"})
}

// Status gauges metrics. Provides information about running systems

func currentCapturers() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "capturer_up",
		Help:      "Capturers that are up",
	}, []string{"uuid", "capturer_uuid", "capturer_url"})
}
