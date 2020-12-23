package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const namespace = "aton"
const subsystem = "ctl"

var (
	capturerReceivedFramesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "capturer_received_frame_total",
		Help:      "The total frame processed by capturers and received in CTL",
	}, []string{"uuid"})

	capturerFailedFramesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "capturer_failed_frame_total",
		Help:      "The total failed frames returned by capturers",
	}, []string{"uuid"})

	processedFramesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "processed_frames_total",
		Help:      "The total frame processed detectors and received in CTL",
	}, []string{"uuid"})

	failedFramesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "unrecognized_frames_total",
		Help:      "The total unrecognized frames processed by detectors and received in CTL",
	}, []string{"uuid"})

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
