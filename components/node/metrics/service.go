package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Service struct {
	registry                    *prometheus.Registry
	capturerReceivedFramesTotal *prometheus.CounterVec
	capturerReceivedFramesBytes *prometheus.HistogramVec
	capturerFailedFramesTotal   *prometheus.CounterVec
	processingFrameTimeSeconds  *prometheus.HistogramVec
	processedFramesTotal        *prometheus.CounterVec
	failedFramesTotal           *prometheus.CounterVec
	unrecognizedEntitiesTotal   *prometheus.CounterVec
	entitiesTotal               *prometheus.CounterVec
	currentCapturers            *prometheus.GaugeVec
}

func NewService() *Service {
	promRegistry := prometheus.NewRegistry()
	s := &Service{
		registry:                    promRegistry,
		capturerReceivedFramesTotal: capturerReceivedFramesTotal(),
		capturerReceivedFramesBytes: capturerReceivedFramesBytes(),
		capturerFailedFramesTotal:   capturerFailedFramesTotal(),
		processingFrameTimeSeconds:  processingFramesDurationSeconds(),
		processedFramesTotal:        processedFramesTotal(),
		failedFramesTotal:           failedFramesTotal(),
		entitiesTotal:               entitiesTotal(),
		unrecognizedEntitiesTotal:   unrecognizedEntitiesTotal(),
		currentCapturers:            currentCapturers(),
	}
	s.registerMetrics(promRegistry)
	return s
}

func (s *Service) registerMetrics(reg *prometheus.Registry) {
	reg.MustRegister(prometheus.NewGoCollector())
	reg.MustRegister(s.capturerReceivedFramesTotal)
	reg.MustRegister(s.capturerReceivedFramesBytes)
	reg.MustRegister(s.capturerFailedFramesTotal)
	reg.MustRegister(s.processingFrameTimeSeconds)
	reg.MustRegister(s.processedFramesTotal)
	reg.MustRegister(s.failedFramesTotal)
	reg.MustRegister(s.entitiesTotal)
	reg.MustRegister(s.unrecognizedEntitiesTotal)
	reg.MustRegister(s.currentCapturers)
}

func (s *Service) IncCapturerReceivedFramesTotal(uuid, url string) {
	s.capturerReceivedFramesTotal.WithLabelValues(uuid, url).Inc()
}

func (s *Service) IncCapturerReceivedFramesBytes(uuid, url string, bytes int) {
	s.capturerReceivedFramesBytes.WithLabelValues(uuid, url).Observe(float64(bytes))
}

func (s *Service) IncCapturerFailedFramesTotal(uuid, url string) {
	s.capturerFailedFramesTotal.WithLabelValues(uuid, url).Inc()
}

func (s *Service) IncProcessedTimeFramesSeconds(seconds float64) {
	s.processingFrameTimeSeconds.WithLabelValues().Observe(seconds)
}

func (s *Service) IncProcessedFramesTotal() {
	s.processedFramesTotal.WithLabelValues().Inc()
}

func (s *Service) IncFailedFramesTotal() {
	s.failedFramesTotal.WithLabelValues().Inc()
}

func (s *Service) AddEntitiesTotal(count int) {
	s.entitiesTotal.WithLabelValues().Add(float64(count))
}

func (s *Service) AddUnrecognizedEntitiesTotal(count int) {
	s.unrecognizedEntitiesTotal.WithLabelValues().Add(float64(count))
}

func (s *Service) CapturerUP(uuid, url string) {
	s.currentCapturers.WithLabelValues(uuid, url).Inc()
}

func (s *Service) CapturerDown(uuid, url string) {
	s.currentCapturers.WithLabelValues(uuid, url).Dec()
}

func (s *Service) MustRegister(c ...prometheus.Collector) {
	s.registry.MustRegister(c...)
}

func (s *Service) HTTPHandler() http.Handler {
	return promhttp.HandlerFor(s.registry, promhttp.HandlerOpts{})
}
