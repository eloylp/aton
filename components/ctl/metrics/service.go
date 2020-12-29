package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Service struct {
	registry                    *prometheus.Registry
	capturerReceivedFramesTotal *prometheus.CounterVec
	capturerFailedFramesTotal   *prometheus.CounterVec
	processedFramesTotal        *prometheus.CounterVec
	failedFramesTotal           *prometheus.CounterVec
	unrecognizedFramesTotal     *prometheus.CounterVec
	currentCapturers            *prometheus.GaugeVec
	currentDetectors            *prometheus.GaugeVec
}

func NewService() *Service {
	promRegistry := prometheus.NewRegistry()
	s := &Service{
		registry:                    promRegistry,
		capturerReceivedFramesTotal: capturerReceivedFramesTotal(),
		capturerFailedFramesTotal:   capturerFailedFramesTotal(),
		processedFramesTotal:        processedFramesTotal(),
		failedFramesTotal:           failedFramesTotal(),
		unrecognizedFramesTotal:     unrecognizedFramesTotal(),
		currentCapturers:            currentCapturers(),
		currentDetectors:            currentDetectors(),
	}
	s.registerMetrics(promRegistry)
	return s
}

func (s *Service) registerMetrics(reg *prometheus.Registry) {
	reg.MustRegister(prometheus.NewGoCollector())
	reg.MustRegister(s.capturerReceivedFramesTotal)
	reg.MustRegister(s.capturerFailedFramesTotal)
	reg.MustRegister(s.processedFramesTotal)
	reg.MustRegister(s.failedFramesTotal)
	reg.MustRegister(s.unrecognizedFramesTotal)
	reg.MustRegister(s.currentCapturers)
	reg.MustRegister(s.currentDetectors)
}

func (s *Service) IncCapturerReceivedFramesTotal(labelValues ...string) {
	s.capturerReceivedFramesTotal.WithLabelValues(labelValues...).Inc()
}

func (s *Service) IncCapturerFailedFramesTotal(labelValues ...string) {
	s.capturerFailedFramesTotal.WithLabelValues(labelValues...).Inc()
}

func (s *Service) IncProcessedFramesTotal(labelValues ...string) {
	s.processedFramesTotal.WithLabelValues(labelValues...).Inc()
}

func (s *Service) IncFailedFramesTotal(labelValues ...string) {
	s.failedFramesTotal.WithLabelValues(labelValues...).Inc()
}

func (s *Service) IncUnrecognizedFramesTotal(labelValues ...string) {
	s.unrecognizedFramesTotal.WithLabelValues(labelValues...).Inc()
}

func (s *Service) CapturerUP(labelValues ...string) {
	s.currentCapturers.WithLabelValues(labelValues...).Inc()
}

func (s *Service) CapturerDown(labelValues ...string) {
	s.currentCapturers.WithLabelValues(labelValues...).Dec()
}

func (s *Service) DetectorUP(labelValues ...string) {
	s.currentDetectors.WithLabelValues(labelValues...).Inc()
}

func (s *Service) DetectorDown(labelValues ...string) {
	s.currentDetectors.WithLabelValues(labelValues...).Dec()
}

func (s *Service) MustRegister(c ...prometheus.Collector) {
	s.registry.MustRegister(c...)
}

func (s *Service) HTTPHandler() http.Handler {
	return promhttp.HandlerFor(s.registry, promhttp.HandlerOpts{})
}
