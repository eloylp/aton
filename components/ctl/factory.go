package ctl

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/eloylp/aton/components/ctl/config"
	"github.com/eloylp/aton/components/ctl/metrics"
	"github.com/eloylp/aton/components/ctl/www"
)

func New(opts ...config.Option) (*Server, error) {
	cfg := &config.Config{}
	for _, o := range opts {
		o(cfg)
	}
	return newWithConfig(cfg)
}

func newWithConfig(cfg *config.Config) (*Server, error) {
	logger := logrus.New()
	if cfg.LogFormat == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
	if cfg.LogOutput == nil {
		cfg.LogOutput = os.Stdout
	}
	logger.SetOutput(cfg.LogOutput)
	metricsService := metrics.NewService()
	c := NewWith(metricsService, logger,
		config.WithListenAddress(cfg.ListenAddress),
	)
	return c, nil
}

func NewFromEnv() (*Server, error) {
	cfg, err := config.FromEnv()
	if err != nil {
		return nil, fmt.Errorf("ctl: %w", err)
	}
	return newWithConfig(cfg)
}

func NewWith(metricsService *metrics.Service, logger *logrus.Logger, opts ...config.Option) *Server {
	cfg := &config.Config{
		APIReadTimeout:  config.DefaultAPIReadTimeout,
		APIWriteTimeout: config.DefaultAPIWriteTimeout,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	api := &http.Server{
		Addr:         cfg.ListenAddress,
		Handler:      www.Router(metricsService.HTTPHandler()),
		ReadTimeout:  cfg.APIReadTimeout,
		WriteTimeout: cfg.APIWriteTimeout,
	}
	ctl := &Server{
		cfg:            cfg,
		metricsService: metricsService,
		api:            api,
		logger:         logger,
		wg:             &sync.WaitGroup{},
		L:              &sync.Mutex{},
	}
	return ctl
}
