package detector

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/eloylp/aton/components/detector/config"
	"github.com/eloylp/aton/components/detector/metrics"
)

func New(opts ...config.Option) (*Server, error) {
	cfg := &config.Config{}
	for _, o := range opts {
		o(cfg)
	}
	return newWithConfig(cfg)
}

func NewFromEnv() (*Server, error) {
	cfg, err := config.FromEnv()
	if err != nil {
		return nil, fmt.Errorf("detector: %w", err)
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
	if cfg.UUID == "" {
		cfg.UUID = uuid.New().String()
	}
	logger.SetOutput(cfg.LogOutput)
	faceDetector, err := NewGoFace(cfg.ModelDir)
	if err != nil {
		return nil, fmt.Errorf("detector: %w", err)
	}
	m := metrics.NewService(cfg.UUID)
	capturerHandler := NewCapturerHandler(logger, m, 100) // todo parametrize
	service := NewService(cfg.UUID, faceDetector, capturerHandler, m, logger, time.Now)
	server := NewServer(cfg.ListenAddr, service, cfg.MetricsListenAddr, m, logger)
	return server, nil
}
