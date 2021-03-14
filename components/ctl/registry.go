package ctl

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

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

func (r *NodeRegistry) LoadCategories(ctx context.Context, categories []string, image []byte) error {
	wg := sync.WaitGroup{}
	opErrors := make(chan error, len(r.registry))
	for _, v := range r.registry {
		wg.Add(1)
		go func(v *NodeHandler) {
			select {
			case <-ctx.Done():
				opErrors <- fmt.Errorf("nodeRegistry: categories programming aborted: %s", ctx.Err())
			default:
				if err := v.LoadCategories(ctx, categories, image); err != nil {
					opErrors <- err
					r.logger.Errorf("nodeRegistry: categories programming error for node %s: %s", v.node.UUID, err)
				}
			}
		}(v)
	}
	wg.Wait()
	close(opErrors)
	message := strings.Builder{}
	if len(opErrors) > 0 {
		for e := range opErrors {
			message.WriteString(e.Error())
		}
		return errors.New("result from programming categories:" + message.String())
	}
	return nil
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
