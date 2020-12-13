package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/eloylp/aton/internal/logging"
	"github.com/eloylp/aton/internal/proto"
)

type DetectorClient struct {
	addr           string
	client         proto.DetectorClient
	internalClient proto.Detector_RecognizeClient
	logger         logging.Logger
}

func NewDetectorClient(addr string, logger logging.Logger) *DetectorClient {
	return &DetectorClient{addr: addr, logger: logger}
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

func (c *DetectorClient) SendToRecognize(req *proto.RecognizeRequest) error {
	if err := c.internalClient.Send(req); err != nil {
		return fmt.Errorf("detectorclient: recognize: send: %w", err)
	}
	return nil
}

func (c *DetectorClient) NextRecognizeResponse() (*proto.RecognizeResponse, error) {
	resp, err := c.internalClient.Recv()
	if err != nil {
		return nil, fmt.Errorf("detectorclient: recognize: next: %w", err)
	}
	return resp, nil
}

func (c *DetectorClient) StartRecognize(ctx context.Context) error {
	var err error
	c.internalClient, err = c.client.Recognize(ctx)
	if err != nil {
		return fmt.Errorf("detectorclient: recognize: %w", err)
	}
	return nil
}

func (c *DetectorClient) Shutdown() error {
	if err := c.internalClient.CloseSend(); err != nil {
		return fmt.Errorf("detectorclient: shutdown: %w", err)
	}
	return nil
}
