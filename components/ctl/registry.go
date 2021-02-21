package ctl

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

type DetectorRegistry struct {
	registry              map[string]*DetectorHandler
	detectorPriorityQueue DetectorPriorityQueue
	logger                *logrus.Logger
}

func NewDetectorRegistry(detectorQueue DetectorPriorityQueue, logger *logrus.Logger) *DetectorRegistry {
	return &DetectorRegistry{
		registry:              make(map[string]*DetectorHandler),
		detectorPriorityQueue: detectorQueue,
		logger:                logger,
	}
}

func (r *DetectorRegistry) Add(s *DetectorHandler) {
	s.priorityQueue = r.detectorPriorityQueue
	r.registry[s.detector.UUID] = s
}

func (r *DetectorRegistry) Find(_ string) (DetectorClient, error) {
	return nil, nil
}

func (r *DetectorRegistry) ShutdownAll(ctx context.Context) error {
	for _, v := range r.registry {
		select {
		case <-ctx.Done():
			return fmt.Errorf("detectorRegistry: shutdown aborted: %s", ctx.Err())
		default:
			if err := v.Shutdown(); err != nil {
				r.logger.Errorf("detectorRegistry: shutdown error for detector %s: %s", v.detector.UUID, err)
			}
		}
	}
	return nil
}
