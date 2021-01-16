package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Service struct {
	UUID                        string
	registry                    *prometheus.Registry
	capturerReceivedFramesTotal *prometheus.CounterVec
	capturerFailedFramesTotal   *prometheus.CounterVec
	processedFramesTotal        *prometheus.CounterVec
	failedFramesTotal           *prometheus.CounterVec
	unrecognizedFramesTotal     *prometheus.CounterVec
	currentCapturers            *prometheus.GaugeVec
}

func NewService(detectorUuid string) *Service {
	promRegistry := prometheus.NewRegistry()
	s := &Service{
		UUID:                        detectorUuid,
		registry:                    promRegistry,
		capturerReceivedFramesTotal: capturerReceivedFramesTotal(),
		capturerFailedFramesTotal:   capturerFailedFramesTotal(),
		processedFramesTotal:        processedFramesTotal(),
		failedFramesTotal:           failedFramesTotal(),
		unrecognizedFramesTotal:     unrecognizedFramesTotal(),
		currentCapturers:            currentCapturers(),
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
}

func (s *Service) IncCapturerReceivedFramesTotal(capturerUuid string) {
	s.capturerReceivedFramesTotal.WithLabelValues(s.UUID, capturerUuid).Inc()
}

func (s *Service) IncCapturerFailedFramesTotal(capturerUuid string) {
	s.capturerFailedFramesTotal.WithLabelValues(s.UUID, capturerUuid).Inc()
}

func (s *Service) IncProcessedFramesTotal() {
	s.processedFramesTotal.WithLabelValues(s.UUID).Inc()
}

func (s *Service) IncFailedFramesTotal() {
	s.failedFramesTotal.WithLabelValues(s.UUID).Inc()
}

func (s *Service) IncUnrecognizedFramesTotal() {
	s.unrecognizedFramesTotal.WithLabelValues(s.UUID).Inc()
}

func (s *Service) CapturerUP(capturerUuid string) {
	s.currentCapturers.WithLabelValues(s.UUID, capturerUuid).Inc()
}

func (s *Service) CapturerDown(capturerUuid string) {
	s.currentCapturers.WithLabelValues(s.UUID, capturerUuid).Dec()
}

func (s *Service) MustRegister(c ...prometheus.Collector) {
	s.registry.MustRegister(c...)
}

func (s *Service) HTTPHandler() http.Handler {
	return promhttp.HandlerFor(s.registry, promhttp.HandlerOpts{})
}
