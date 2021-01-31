package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Service struct {
	registry         *prometheus.Registry
	currentDetectors *prometheus.GaugeVec
}

func NewService() *Service {
	promRegistry := prometheus.NewRegistry()
	s := &Service{
		registry:         promRegistry,
		currentDetectors: currentDetectors(),
	}
	s.registerMetrics(promRegistry)
	return s
}

func (s *Service) registerMetrics(reg *prometheus.Registry) {
	reg.MustRegister(prometheus.NewGoCollector())
	reg.MustRegister(s.currentDetectors)
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
