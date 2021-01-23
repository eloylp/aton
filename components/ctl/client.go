package ctl

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/eloylp/aton/components/ctl/metrics"
	"github.com/eloylp/aton/components/proto"
)

const (
	MaxDetectorStatusQueueSize = 50
	MaxDetectorResultQueueSize = 50
	StreamConnectTimeout       = time.Second
)

type GRPCDetectorClient struct {
	addr                 string
	client               proto.DetectorClient
	clientConn           *grpc.ClientConn
	logger               *logrus.Logger
	metricsRegistry      *metrics.Service
	detectorStatusQueue  chan *Status
	detectorResultsQueue chan *Result
	shutdown             chan struct{}
	wg                   *sync.WaitGroup
}

func NewGRPCDetectorClient(addr string, logger *logrus.Logger, metricsRegistry *metrics.Service) *GRPCDetectorClient {
	return &GRPCDetectorClient{
		addr:                 addr,
		logger:               logger,
		metricsRegistry:      metricsRegistry,
		detectorStatusQueue:  make(chan *Status, MaxDetectorStatusQueueSize),
		detectorResultsQueue: make(chan *Result, MaxDetectorResultQueueSize),
		shutdown:             make(chan struct{}),
		wg:                   &sync.WaitGroup{},
	}
}

func (c *GRPCDetectorClient) Connect() error {
	clientMetrics := grpc_prometheus.NewClientMetrics()
	clientMetrics.EnableClientHandlingTimeHistogram()
	c.metricsRegistry.MustRegister(clientMetrics)
	logrusEntry := logrus.NewEntry(c.logger)
	backOffScalar := 500 * time.Millisecond
	backOffJitter := 0.35
	var err error
	c.clientConn, err = grpc.Dial(c.addr,
		grpc.WithInsecure(),
		grpc.WithStreamInterceptor(
			grpc_middleware.ChainStreamClient(
				clientMetrics.StreamClientInterceptor(),
				grpc_logrus.StreamClientInterceptor(logrusEntry),
				grpc_retry.StreamClientInterceptor(grpc_retry.WithBackoff(
					grpc_retry.BackoffExponentialWithJitter(backOffScalar, backOffJitter),
				)),
			),
		),
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				clientMetrics.UnaryClientInterceptor(),
				grpc_logrus.UnaryClientInterceptor(logrusEntry),
				grpc_retry.UnaryClientInterceptor(grpc_retry.WithBackoff(
					grpc_retry.BackoffExponentialWithJitter(backOffScalar, backOffJitter),
				)),
			),
		),
	)
	if err != nil {
		return fmt.Errorf("dectectorclient: %w", err)
	}
	c.client = proto.NewDetectorClient(c.clientConn)
	return nil
}

func (c *GRPCDetectorClient) LoadCategories(ctx context.Context, request *LoadCategoriesRequest) error {
	if _, err := c.client.LoadCategories(ctx, &proto.LoadCategoriesRequest{
		Categories: request.Categories,
		Image:      request.Image,
	}); err != nil {
		return err
	}
	return nil
}

func (c *GRPCDetectorClient) AddCapturer(ctx context.Context, request *AddCapturerRequest) error {
	if _, err := c.client.AddCapturer(ctx, &proto.AddCapturerRequest{
		CapturerUuid: request.UUID,
		CapturerUrl:  request.URL,
	}); err != nil {
		return err
	}
	return nil
}

func (c *GRPCDetectorClient) RemoveCapturer(ctx context.Context, request *RemoveCapturerRequest) error {
	if _, err := c.client.RemoveCapturer(ctx, &proto.RemoveCapturerRequest{
		CapturerUuid: request.UUID,
	}); err != nil {
		return err
	}
	return nil
}

func (c *GRPCDetectorClient) startStatusProc(interval time.Duration) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		ctx, cancl := context.WithTimeout(context.Background(), StreamConnectTimeout)
		defer cancl()
		stream, err := c.client.InformStatus(ctx, &proto.InformStatusRequest{
			Interval: durationpb.New(interval),
		})
		if err != nil {
			c.logger.Errorf("gRPCDetectorClient: %v", err)
			return
		}
	mainLoop:
		for {
			status, err := stream.Recv()
			if err == io.EOF {
				c.logger.Info("gRPCDetectorClient: result stream ended.")
				return
			}
			if err != nil {
				c.logger.Errorf("gRPCDetectorClient: %v", err)
				continue
			}
			select {
			case <-c.shutdown:
				break mainLoop
			default:
				var capturers []*Capturer
				for _, c := range status.Capturers {
					capturers = append(capturers, &Capturer{
						UUID:   c.Uuid,
						URL:    c.Url,
						Status: c.Status.String(),
					})
				}
				c.detectorStatusQueue <- &Status{
					Description: status.Description,
					Capturers:   capturers,
				}
			}
		}
	}()
}

func (c *GRPCDetectorClient) NextStatus() (*Status, error) {
	status, ok := <-c.detectorStatusQueue
	if !ok {
		return nil, io.EOF
	}
	return status, nil
}

func (c *GRPCDetectorClient) startResultsProc() {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		ctx, cancl := context.WithTimeout(context.Background(), StreamConnectTimeout)
		defer cancl()
		stream, err := c.client.ProcessResults(ctx, &empty.Empty{})
		if err != nil {
			c.logger.Errorf("gRPCDetectorClient: %v", err)
			return
		}
	mainLoop:
		for {
			result, err := stream.Recv()
			if err != io.EOF {
				c.logger.Info("gRPCDetectorClient: result stream ended.")
				return
			}
			if err != nil {
				c.logger.Errorf("gRPCDetectorClient: %v", err)
				continue
			}
			select {
			case <-c.shutdown:
				break mainLoop
			default:
				c.detectorResultsQueue <- &Result{
					DetectorUUID:  result.DetectorUuid,
					Recognized:    result.Recognized,
					TotalEntities: result.TotalEntities,
					RecognizedAt:  result.RecognizedAt.AsTime(),
					CapturedAt:    result.RecognizedAt.AsTime(),
				}
			}
		}
	}()
}

func (c *GRPCDetectorClient) NextResult() (*Result, error) {
	result, ok := <-c.detectorResultsQueue
	if !ok {
		return nil, io.EOF
	}
	return result, nil
}

func (c *GRPCDetectorClient) Shutdown() error {
	close(c.shutdown)
	c.wg.Wait()
	if err := c.clientConn.Close(); err != nil {
		return fmt.Errorf("detectorclient: shutdown: %w", err)
	}
	return nil
}

type LoadCategoriesRequest struct {
	Categories []string
	Image      []byte
}

type AddCapturerRequest struct {
	UUID string
	URL  string
}

type RemoveCapturerRequest struct {
	UUID string
}

type Status struct {
	Description string
	Capturers   []*Capturer
}

type Capturer struct {
	UUID   string
	URL    string
	Status string
}

type Result struct {
	DetectorUUID  string
	Recognized    []string
	TotalEntities int32
	RecognizedAt  time.Time
	CapturedAt    time.Time
}
