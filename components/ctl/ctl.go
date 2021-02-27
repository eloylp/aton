package ctl

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/eloylp/aton/components/ctl/metrics"
)

type ClientGenerator func(addr string, l *logrus.Logger, service *metrics.Service) NodeClient

var (
	ErrNotAvailableNode = errors.New("ctl: cannot find suitable node")
)

type Ctl struct {
	registry        *NodeRegistry
	nodeQueue       NodePriorityQueue
	logger          *logrus.Logger
	metricsService  *metrics.Service
	clientGenerator ClientGenerator
}

func NewCtl(logger *logrus.Logger, metricsService *metrics.Service, clientGen ClientGenerator) *Ctl {
	queue := NewHeapNodePriorityQueue()
	return &Ctl{
		registry:        NewNodeRegistry(queue, logger),
		nodeQueue:       queue,
		logger:          logger,
		metricsService:  metricsService,
		clientGenerator: clientGen,
	}
}

func (c *Ctl) AddNode(addr string) (string, error) {
	client := c.clientGenerator(addr, c.logger, c.metricsService)
	uid := uuid.New().String()
	node := &Node{UUID: uid, Addr: addr}
	nodeHandler := NewNodeHandler(node, client, c.logger)
	c.registry.Add(nodeHandler)
	if err := nodeHandler.Start(); err != nil {
		return "", fmt.Errorf("ctl: error adding node: %w", err)
	}
	c.logger.Infof("ctl: added node at %s for %s", addr, uid)
	return uid, nil
}

func (c *Ctl) AddCapturer(ctx context.Context, uid, url string) error {
	electedNode := c.nodeQueue.Next()
	if electedNode == nil {
		c.logger.Errorf("ctl: not suitable node for capturer %s with URL %s", uid, url)
		return ErrNotAvailableNode
	}
	client, err := c.registry.Find(electedNode.UUID)
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
