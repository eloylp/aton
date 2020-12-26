package detector

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/eloylp/aton/internal/detector/config"
	"github.com/eloylp/aton/internal/detector/engine"
	"github.com/eloylp/aton/internal/detector/grpc"
)

func New(opts ...config.Option) (*grpc.Server, error) {
	cfg := &config.Config{}
	for _, o := range opts {
		o(cfg)
	}
	return NewWithConfig(cfg)
}

func NewFromEnv() (*grpc.Server, error) {
	cfg, err := config.FromEnv()
	if err != nil {
		return nil, fmt.Errorf("detector: %w", err)
	}
	return NewWithConfig(cfg)
}

func NewWithConfig(cfg *config.Config) (*grpc.Server, error) {
	logger := logrus.New()
	if cfg.LogFormat == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
	faceDetector, err := engine.NewGoFaceDetector(cfg.ModelDir)
	if err != nil {
		return nil, fmt.Errorf("detector: %w", err)
	}
	service := grpc.NewService(cfg.UUID, faceDetector, logger, time.Now)
	server := grpc.NewServer(cfg.ListenAddr, service, cfg.MetricsListenAddr, logger)
	return server, nil
}
