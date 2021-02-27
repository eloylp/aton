package ctl

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

type NodeRegistry struct {
	registry          map[string]*NodeHandler
	nodePriorityQueue NodePriorityQueue
	logger            *logrus.Logger
}

func NewNodeRegistry(nodeQueue NodePriorityQueue, logger *logrus.Logger) *NodeRegistry {
	return &NodeRegistry{
		registry:          make(map[string]*NodeHandler),
		nodePriorityQueue: nodeQueue,
		logger:            logger,
	}
}

func (r *NodeRegistry) Add(s *NodeHandler) {
	s.priorityQueue = r.nodePriorityQueue
	r.registry[s.node.UUID] = s
}

func (r *NodeRegistry) Find(_ string) (NodeClient, error) {
	return nil, nil
}

func (r *NodeRegistry) ShutdownAll(ctx context.Context) error {
	for _, v := range r.registry {
		select {
		case <-ctx.Done():
			return fmt.Errorf("nodeRegistry: shutdown aborted: %s", ctx.Err())
		default:
			if err := v.Shutdown(); err != nil {
				r.logger.Errorf("nodeRegistry: shutdown error for node %s: %s", v.node.UUID, err)
			}
		}
	}
	return nil
}
