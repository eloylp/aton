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
	}, []string{"uuid", "url"})
}

func capturerReceivedFramesBytes() *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "capturer_received_frames_bytes",
		Help:      "The size of frames obtained by a capturer",
	}, []string{"uuid", "url"})
}

func capturerFailedFramesTotal() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "capturer_failed_frames_total",
		Help:      "The failed frames returned by capturers",
	}, []string{"uuid", "url"})
}

// Responses of nodes metrics

func processingFramesDurationSeconds() *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "frame_processing_time_seconds",
		Help:      "The time a node spends processing frames",
	}, []string{})
}

func processedFramesTotal() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "processed_frames_total",
		Help:      "The total frame processed by this node",
	}, []string{})
}

func failedFramesTotal() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "error_frames_total",
		Help:      "The failed, not processed, frames on this node",
	}, []string{})
}

func entitiesTotal() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "entities_total",
		Help:      "The total entities in frames processed by nodes",
	}, []string{})
}

func unrecognizedEntitiesTotal() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "unrecognized_entities_total",
		Help:      "Entities in frames that are not classified by ML models",
	}, []string{})
}

// Status gauges metrics. Provides information about running systems

func currentCapturers() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "capturer_up",
		Help:      "Capturers that are up on this node",
	}, []string{"uuid", "url"})
}
