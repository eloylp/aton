package detector

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

func New(opts ...Option) (*Server, error) {
	cfg := &Config{}
	for _, o := range opts {
		o(cfg)
	}
	return NewWithConfig(cfg)
}

func NewFromEnv() (*Server, error) {
	cfg, err := FromEnv()
	if err != nil {
		return nil, fmt.Errorf("detector: %w", err)
	}
	return NewWithConfig(cfg)
}

func NewWithConfig(cfg *Config) (*Server, error) {
	logger := logrus.New()
	if cfg.LogFormat == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
	faceDetector, err := NewGoFace(cfg.ModelDir)
	if err != nil {
		return nil, fmt.Errorf("detector: %w", err)
	}
	service := NewService(cfg.UUID, faceDetector, logger, time.Now)
	server := NewServer(cfg.ListenAddr, service, cfg.MetricsListenAddr, logger)
	return server, nil
}
