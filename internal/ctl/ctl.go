package ctl

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/eloylp/aton/internal/ctl/config"
	"github.com/eloylp/aton/internal/ctl/engine"
	"github.com/eloylp/aton/internal/ctl/grpc"
	"github.com/eloylp/aton/internal/ctl/metrics"
)

func New(opts ...config.Option) (*engine.Ctl, error) {
	cfg := &config.Config{}
	for _, o := range opts {
		o(cfg)
	}
	return NewWithConfig(cfg)
}

func NewWithConfig(cfg *config.Config) (*engine.Ctl, error) {
	logger := logrus.New()
	if cfg.LogFormat == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
	metricsService := metrics.NewService()
	dc := grpc.NewDetectorClient(cfg.Detector.Address, logger, metricsService)
	c := engine.New(dc, metricsService, logger,
		config.WithListenAddress(cfg.ListenAddress),
	)
	return c, nil
}

func NewFromEnv() (*engine.Ctl, error) {
	cfg, err := config.FromEnv()
	if err != nil {
		return nil, fmt.Errorf("ctl: %w", err)
	}
	return NewWithConfig(cfg)
}
