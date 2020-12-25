package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const namespace = "aton"
const subsystem = "ctl"

var (
	// Capturer (input video stream) metrics
	capturerReceivedFramesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "capturer_received_frames_total",
		Help:      "The total frames obtained from capturers",
	}, []string{"uuid"})

	capturerFailedFramesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "capturer_failed_frames_total",
		Help:      "The failed frames returned by capturers",
	}, []string{"uuid"})

	// Responses of detectors metrics
	processedFramesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "processed_frames_total",
		Help:      "The total frame processed detectors",
	}, []string{"uuid"})

	failedFramesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "error_frames_total",
		Help:      "The failed, not processed, frames received from detectors",
	}, []string{"uuid"})

	unrecognizedFramesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "unrecognized_frames_total",
		Help:      "The total unrecognized frames processed by detectors",
	}, []string{"uuid"})

	// Status gauges metrics. Provides information about running systems
	currentCapturers = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "capturer_up",
		Help:      "Capturers that are up",
	}, []string{"uuid"})

	currentDetectors = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "detector_up",
		Help:      "Detectors that are up",
	}, []string{"uuid"})
)
