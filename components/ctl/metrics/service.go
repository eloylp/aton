package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Service struct {
	registry     *prometheus.Registry
	currentNodes *prometheus.GaugeVec
}

func NewService() *Service {
	promRegistry := prometheus.NewRegistry()
	s := &Service{
		registry:     promRegistry,
		currentNodes: currentNodes(),
	}
	s.registerMetrics(promRegistry)
	return s
}

func (s *Service) registerMetrics(reg *prometheus.Registry) {
	reg.MustRegister(prometheus.NewGoCollector())
	reg.MustRegister(s.currentNodes)
}

func (s *Service) NodeUP(labelValues ...string) {
	s.currentNodes.WithLabelValues(labelValues...).Inc()
}

func (s *Service) NodeDown(labelValues ...string) {
	s.currentNodes.WithLabelValues(labelValues...).Dec()
}

func (s *Service) MustRegister(c ...prometheus.Collector) {
	s.registry.MustRegister(c...)
}

func (s *Service) HTTPHandler() http.Handler {
	return promhttp.HandlerFor(s.registry, promhttp.HandlerOpts{})
}
