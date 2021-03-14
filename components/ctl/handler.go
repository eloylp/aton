package ctl

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

type NodeHandler struct {
	node           *Node
	client         NodeClient
	priorityQueue  NodePriorityQueue
	processStopper chan struct{}
	wg             *sync.WaitGroup
	logger         *logrus.Logger
}

func NewNodeHandler(node *Node, client NodeClient, logger *logrus.Logger) *NodeHandler {
	return &NodeHandler{
		node:           node,
		client:         client,
		logger:         logger,
		wg:             &sync.WaitGroup{},
		processStopper: make(chan struct{}),
	}
}

func (ds *NodeHandler) Start() error {
	if err := ds.client.Connect(); err != nil {
		return fmt.Errorf("ctl: could not connect to node %s: %w", ds.node.UUID, err)
	}
	ds.wg.Add(2)
	go ds.processStatus()
	go ds.processResults()
	return nil
}

func (ds *NodeHandler) LoadCategories(ctx context.Context, categories []string, image []byte) error {
	if err := ds.client.LoadCategories(ctx, &LoadCategoriesRequest{
		Categories: categories,
		Image:      image,
	}); err != nil {
		return fmt.Errorf("ctl: could not load categories in node %s: %w", ds.node.UUID, err)
	}
	return nil
}

func (ds *NodeHandler) processResults() {
	defer ds.wg.Done()
	for {
		select {
		case <-ds.processStopper:
			ds.logger.Infof("nodeHandler: closed processing results of %s", ds.node.UUID)
			return
		default:
			r, err := ds.client.NextResult()
			if err != nil {
				ds.logger.Errorf("nodeHandler: error obtaining next result: %s", err)
				return
			}
			resultFormat := "ctl: result: %s - %d (%s) - %d | %s | %s"
			ds.logger.Infof(resultFormat,
				r.CapturerUUID,
				len(r.Recognized),
				strings.Join(r.Recognized, ","),
				r.TotalEntities,
				r.CapturedAt,
				r.RecognizedAt)
		}
	}
}

func (ds *NodeHandler) processStatus() {
	defer ds.wg.Done()
	for {
		select {
		case <-ds.processStopper:
			ds.logger.Infof("nodeHandler: closed processing status of %s", ds.node.UUID)
			return
		default:
			s, err := ds.client.NextStatus()
			if err != nil {
				ds.logger.Errorf("nodeHandler: error obtaining next status: %s", err)
				return
			}
			ds.node.Status = s
			ds.priorityQueue.Upsert(ds.node)
		}
	}
}

func (ds *NodeHandler) Shutdown() error {
	close(ds.processStopper)
	ds.wg.Wait()
	return nil
}
