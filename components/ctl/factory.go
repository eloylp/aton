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

func New(opts ...config.Option) (*Ctl, error) {
	cfg := &config.Config{}
	for _, o := range opts {
		o(cfg)
	}
	return NewWithConfig(cfg)
}

func NewWithConfig(cfg *config.Config) (*Ctl, error) {
	logger := logrus.New()
	if cfg.LogFormat == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
	metricsService := metrics.NewService()
	dc := NewGRPCDetectorClient(cfg.Detector.Address, logger, metricsService)
	c := NewWith(dc, metricsService, logger,
		config.WithListenAddress(cfg.ListenAddress),
	)
	return c, nil
}

func NewFromEnv() (*Ctl, error) {
	cfg, err := config.FromEnv()
	if err != nil {
		return nil, fmt.Errorf("ctl: %w", err)
	}
	return NewWithConfig(cfg)
}

func NewWith(dc DetectorClient, metricsService *metrics.Service, logger *logrus.Logger, opts ...config.Option) *Ctl {
	cfg := &config.Config{
		APIReadTimeout:  config.DefaultAPIReadTimeout,
		APIWriteTimeout: config.DefaultAPIWriteTimeout,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.LogOutput == nil {
		cfg.LogOutput = os.Stdout
	}
	logger.SetOutput(cfg.LogOutput)
	api := &http.Server{
		Addr:         cfg.ListenAddress,
		Handler:      www.Router(metricsService.HTTPHandler()),
		ReadTimeout:  cfg.APIReadTimeout,
		WriteTimeout: cfg.APIWriteTimeout,
	}
	metricsService.DetectorUP(cfg.Detector.UUID)
	ctl := &Ctl{
		cfg:            cfg,
		detectorClient: dc,
		metricsService: metricsService,
		capturers:      CapturerRegistry{},
		api:            api,
		logger:         logger,
		wg:             &sync.WaitGroup{},
		L:              &sync.Mutex{},
	}
	return ctl
}
