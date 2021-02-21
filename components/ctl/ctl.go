package ctl

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/eloylp/aton/components/ctl/metrics"
)

type ClientGenerator func(addr string, l *logrus.Logger, service *metrics.Service) DetectorClient

var (
	ErrNotAvailableDetector = errors.New("ctl: cannot find suitable detector")
)

type Ctl struct {
	registry        *DetectorRegistry
	detectorQueue   DetectorPriorityQueue
	logger          *logrus.Logger
	metricsService  *metrics.Service
	clientGenerator ClientGenerator
}

func NewCtl(logger *logrus.Logger, metricsService *metrics.Service, clientGen ClientGenerator) *Ctl {
	queue := NewHeapDetectorPriorityQueue()
	return &Ctl{
		registry:        NewDetectorRegistry(queue, logger),
		detectorQueue:   queue,
		logger:          logger,
		metricsService:  metricsService,
		clientGenerator: clientGen,
	}
}

func (c *Ctl) AddDetector(addr string) (string, error) {
	client := c.clientGenerator(addr, c.logger, c.metricsService)
	uid := uuid.New().String()
	detector := &Detector{UUID: uid, Addr: addr}
	detectorHandler := NewDetectorHandler(detector, client, c.logger)
	c.registry.Add(detectorHandler)
	if err := detectorHandler.Start(); err != nil {
		return "", fmt.Errorf("ctl: error adding detector: %w", err)
	}
	c.logger.Infof("ctl: added detector at %s for %s", addr, uid)
	return uid, nil
}

func (c *Ctl) AddCapturer(ctx context.Context, uid, url string) error {
	electedDetector := c.detectorQueue.Next()
	if electedDetector == nil {
		c.logger.Errorf("ctl: not suitable node for capturer %s with URL %s", uid, url)
		return ErrNotAvailableDetector
	}
	client, err := c.registry.Find(electedDetector.UUID)
	if err != nil {
		return err
	}
	if err := client.AddCapturer(ctx, &AddCapturerRequest{
		UUID: uid,
		URL:  url,
	}); err != nil {
		return err
	}
	return nil
}

func (c *Ctl) Shutdown(ctx context.Context) error {
	if err := c.registry.ShutdownAll(ctx); err != nil {
		return fmt.Errorf("ctl: error shutting down: %w", err)
	}
	return nil
}
