package grpc

import (
	"context"
	"fmt"
	"io"
	"sync"

	"google.golang.org/grpc"

	"github.com/eloylp/aton/internal/logging"
	"github.com/eloylp/aton/internal/proto"
)

type DetectorClient struct {
	addr            string
	client          proto.DetectorClient
	wg              *sync.WaitGroup
	shutdown        chan struct{}
	detectionInput  chan *proto.RecognizeRequest
	detectionOutput chan *proto.RecognizeResponse
	logger          logging.Logger
}

func NewDetectorClient(addr string, wg *sync.WaitGroup, logger logging.Logger) *DetectorClient {
	return &DetectorClient{addr: addr, wg: wg, logger: logger}
}

func (c *DetectorClient) Connect() error {
	grpcClientConn, err := grpc.Dial(c.addr, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("dectectorclient: %w", err)
	}
	c.client = proto.NewDetectorClient(grpcClientConn)
	return nil
}

func (c *DetectorClient) LoadCategories(ctx context.Context, request *proto.LoadCategoriesRequest) (*proto.LoadCategoriesResponse, error) {
	return c.client.LoadCategories(ctx, request)
}

func (c *DetectorClient) Recognize(ctx context.Context) (chan<- *proto.RecognizeRequest, <-chan *proto.RecognizeResponse, error) {
	gClient, err := c.client.Recognize(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("detectorclient: recognize: %w", err)
	}
	c.wg.Add(2)
	go func() {
		defer func() {
			close(c.detectionOutput)
			c.wg.Done()
		}()
	main:
		for {
			select {
			case <-c.shutdown:
				if err := gClient.CloseSend(); err != nil {
					c.logger.Error(fmt.Errorf("detectorclient: closing error: %w", err))
				}
				close(c.detectionOutput)
				break main
			case request := <-c.detectionInput:
				if err := gClient.Send(request); err != nil {
					c.logger.Error(fmt.Errorf("detectorclient: requesting: %w", err))
				}
			}
		}
	}()
	go func() {
		defer func() {
			close(c.detectionOutput)
			c.wg.Done()
		}()
		for {
			resp, err := gClient.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				c.logger.Error(fmt.Errorf("detectorclient: response: %w", err))
			}
			c.detectionOutput <- resp
			if _, ok := <-c.shutdown; !ok {
				break
			}
		}
	}()
	return c.detectionInput, c.detectionOutput, nil
}

func (c *DetectorClient) Shutdown() {
	c.shutdown <- struct{}{}
	close(c.shutdown)
}
