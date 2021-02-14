package ctl

import (
	"context"
	"net/http"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/eloylp/aton/components/ctl/config"
	"github.com/eloylp/aton/components/ctl/metrics"
)

type Server struct {
	cfg            *config.Config
	metricsService *metrics.Service
	api            *http.Server
	logger         *logrus.Logger
	wg             *sync.WaitGroup
	L              *sync.Mutex
}

func (s *Server) Start() error {
	s.logger.Infof("starting CTL at %s", s.cfg.ListenAddress)
	s.initializeAPI()
	s.wg.Wait()
	return nil
}

func (s *Server) initializeAPI() {
	s.wg.Add(1)
	go func() {
		if err := s.api.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Errorf("ctl: %v", err)
		}
		s.wg.Done()
	}()
}

func (s *Server) Shutdown() {
	s.logger.Info("started graceful shutdown sequence")
	// Close api server
	if err := s.api.Shutdown(context.TODO()); err != nil {
		s.logger.Errorf("ctl: shutdown: %v", err)
	}
	s.wg.Wait()
	s.logger.Infof("stopped CTL at %s", s.cfg.ListenAddress)
}
