package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "aton"
	subsystem = "node"
)

// Capturer (input video stream) metrics

func capturerReceivedFramesTotal() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "capturer_received_frames_total",
		Help:      "The total frames obtained by a capturer",
	}, []string{"uuid", "capturer_uuid", "capturer_url"})
}

func capturerReceivedFramesBytes() *prometheus.SummaryVec {
	return prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "capturer_received_frames_bytes",
		Help:      "The size of frames obtained by a capturer",
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

// Responses of nodes metrics

func processedFramesTotal() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "processed_frames_total",
		Help:      "The total frame processed by this node",
	}, []string{"uuid"})
}

func failedFramesTotal() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "error_frames_total",
		Help:      "The failed, not processed, frames on this node",
	}, []string{"uuid"})
}

func entitiesTotal() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "entities_total",
		Help:      "The total entities in frames processed by nodes",
	}, []string{"uuid"})
}

func unrecognizedEntitiesTotal() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "unrecognized_entities_total",
		Help:      "Entities in frames that are not classified by ML models",
	}, []string{"uuid"})
}

// Status gauges metrics. Provides information about running systems

func currentCapturers() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "capturer_up",
		Help:      "Capturers that are up on this node",
	}, []string{"uuid", "capturer_uuid", "capturer_url"})
}
