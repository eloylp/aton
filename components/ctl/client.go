package ctl

import (
	"context"
	"fmt"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/eloylp/aton/components/ctl/metrics"
	"github.com/eloylp/aton/components/proto"
)

type GRPCDetectorClient struct {
	addr            string
	client          proto.DetectorClient
	internalClient  proto.Detector_RecognizeClient
	clientConn      *grpc.ClientConn
	logger          *logrus.Logger
	metricsRegistry *metrics.Service
}

func NewGRPCDetectorClient(addr string, logger *logrus.Logger, metricsRegistry *metrics.Service) *GRPCDetectorClient {
	return &GRPCDetectorClient{
		addr:            addr,
		logger:          logger,
		metricsRegistry: metricsRegistry,
	}
}

func (c *GRPCDetectorClient) Connect() error {
	clientMetrics := grpc_prometheus.NewClientMetrics()
	clientMetrics.EnableClientHandlingTimeHistogram()
	c.metricsRegistry.MustRegister(clientMetrics)
	logrusEntry := logrus.NewEntry(c.logger)
	var retries uint = 10
	grpcClientConn, err := grpc.Dial(c.addr,
		grpc.WithInsecure(),
		grpc.WithStreamInterceptor(
			grpc_middleware.ChainStreamClient(
				clientMetrics.StreamClientInterceptor(),
				grpc_logrus.StreamClientInterceptor(logrusEntry),
			),
		),
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				clientMetrics.UnaryClientInterceptor(),
				grpc_logrus.UnaryClientInterceptor(logrusEntry),
				grpc_retry.UnaryClientInterceptor(grpc_retry.WithMax(retries)),
			),
		),
	)
	if err != nil {
		return fmt.Errorf("dectectorclient: %w", err)
	}
	c.client = proto.NewDetectorClient(c.clientConn)
	return nil
}

func (c *GRPCDetectorClient) LoadCategories(
	ctx context.Context,
	request *proto.LoadCategoriesRequest) (*proto.LoadCategoriesResponse, error) {
	return c.client.LoadCategories(ctx, request)
}

func (c *GRPCDetectorClient) SendToRecognize(req *proto.RecognizeRequest) error {
	if err := c.internalClient.Send(req); err != nil {
		return fmt.Errorf("detectorclient: recognize: send: %w", err)
	}
	return nil
}

func (c *GRPCDetectorClient) NextRecognizeResponse() (*proto.RecognizeResponse, error) {
	resp, err := c.internalClient.Recv()
	if err != nil {
		return nil, fmt.Errorf("detectorclient: recognize: next: %w", err)
	}
	return resp, nil
}

func (c *GRPCDetectorClient) StartRecognize(ctx context.Context) error {
	var err error
	c.internalClient, err = c.client.Recognize(ctx)
	if err != nil {
		return fmt.Errorf("detectorclient: recognize: %w", err)
	}
	return nil
}

func (c *GRPCDetectorClient) Shutdown() error {
	if c.internalClient != nil {
		if err := c.internalClient.CloseSend(); err != nil {
			return fmt.Errorf("detectorclient: shutdown: %w", err)
		}
	}
	if err := c.clientConn.Close(); err != nil {
		return fmt.Errorf("detectorclient: shutdown: %w", err)
	}
	return nil
}
