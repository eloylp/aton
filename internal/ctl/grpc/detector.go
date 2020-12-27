package grpc

import (
	"context"
	"fmt"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/eloylp/aton/internal/ctl/metrics"
	"github.com/eloylp/aton/internal/proto"
)

type DetectorClient struct {
	addr            string
	client          proto.DetectorClient
	internalClient  proto.Detector_RecognizeClient
	logger          *logrus.Logger
	metricsRegistry *metrics.Service
}

func NewDetectorClient(addr string, logger *logrus.Logger, metricsRegistry *metrics.Service) *DetectorClient {
	return &DetectorClient{
		addr:            addr,
		logger:          logger,
		metricsRegistry: metricsRegistry,
	}
}

func (c *DetectorClient) Connect() error {
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
