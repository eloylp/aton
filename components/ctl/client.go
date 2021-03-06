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
	MaxNodeStatusQueueSize = 50
	MaxNodeResultQueueSize = 50
	StreamConnectTimeout   = time.Second
)

type GRPCNodeClient struct {
	addr            string
	client          proto.NodeClient
	clientConn      *grpc.ClientConn
	logger          *logrus.Logger
	metricsRegistry *metrics.Service
	statusQueue     chan *Status
	resultsQueue    chan *Result
	shutdown        chan struct{}
	wg              *sync.WaitGroup
}

func NewGRPCNodeClient(addr string, logger *logrus.Logger, metricsRegistry *metrics.Service) *GRPCNodeClient {
	return &GRPCNodeClient{
		addr:            addr,
		logger:          logger,
		metricsRegistry: metricsRegistry,
		statusQueue:     make(chan *Status, MaxNodeStatusQueueSize),
		resultsQueue:    make(chan *Result, MaxNodeResultQueueSize),
		shutdown:        make(chan struct{}),
		wg:              &sync.WaitGroup{},
	}
}

const (
	backOffScalar = 500 * time.Millisecond
	backOffJitter = 0.35
)

func (c *GRPCNodeClient) Connect() error {
	grpcMetrics := grpc_prometheus.NewClientMetrics()
	grpcMetrics.EnableClientHandlingTimeHistogram()
	c.metricsRegistry.MustRegister(grpcMetrics)
	logrusEntry := logrus.NewEntry(c.logger)
	var err error
	c.clientConn, err = grpc.Dial(c.addr,
		grpc.WithInsecure(),
		grpc.WithStreamInterceptor(
			grpc_middleware.ChainStreamClient(
				grpcMetrics.StreamClientInterceptor(),
				grpc_logrus.StreamClientInterceptor(logrusEntry),
				grpc_retry.StreamClientInterceptor(grpc_retry.WithBackoff(
					grpc_retry.BackoffExponentialWithJitter(backOffScalar, backOffJitter),
				)),
			),
		),
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				grpcMetrics.UnaryClientInterceptor(),
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
	c.client = proto.NewNodeClient(c.clientConn)
	return nil
}

func (c *GRPCNodeClient) LoadCategories(ctx context.Context, request *LoadCategoriesRequest) error {
	if _, err := c.client.LoadCategories(ctx, &proto.LoadCategoriesRequest{
		Categories: request.Categories,
		Image:      request.Image,
	}); err != nil {
		return err
	}
	return nil
}

func (c *GRPCNodeClient) AddCapturer(ctx context.Context, request *AddCapturerRequest) error {
	if _, err := c.client.AddCapturer(ctx, &proto.AddCapturerRequest{
		CapturerUuid: request.UUID,
		CapturerUrl:  request.URL,
	}); err != nil {
		return err
	}
	return nil
}

func (c *GRPCNodeClient) RemoveCapturer(ctx context.Context, request *RemoveCapturerRequest) error {
	if _, err := c.client.RemoveCapturer(ctx, &proto.RemoveCapturerRequest{
		CapturerUuid: request.UUID,
	}); err != nil {
		return err
	}
	return nil
}

func (c *GRPCNodeClient) startStatusProc(interval time.Duration) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		ctx, cancl := context.WithTimeout(context.Background(), StreamConnectTimeout)
		defer cancl()
		stream, err := c.client.InformStatus(ctx, &proto.InformStatusRequest{
			Interval: durationpb.New(interval),
		})
		if err != nil {
			c.logger.Errorf("gRPCNodeClient: %v", err)
			return
		}
	mainLoop:
		for {
			status, err := stream.Recv()
			if err == io.EOF {
				close(c.statusQueue)
				c.logger.Info("gRPCNodeClient: result stream ended.")
				return
			}
			if err != nil {
				c.logger.Errorf("gRPCNodeClient: %v", err)
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
				c.statusQueue <- c.mapStatus(status, capturers)
			}
		}
	}()
}

func (c *GRPCNodeClient) mapStatus(status *proto.Status, capturers []*Capturer) *Status {
	return &Status{
		Description: status.Description,
		Capturers:   capturers,
		System: System{
			CPUCount: int(status.System.CpuCount),
			Network: Network{
				RxBytesSec: status.System.Network.RxBytesSec,
				TxBytesSec: status.System.Network.TxBytesSec,
			},
			LoadAverage: LoadAverage{
				Avg1:  status.System.LoadAverage.Avg_1,
				Avg5:  status.System.LoadAverage.Avg_5,
				Avg15: status.System.LoadAverage.Avg_15,
			},
			Memory: Memory{
				UsedMemoryBytes:  status.System.Memory.UsedMemoryBytes,
				TotalMemoryBytes: status.System.Memory.TotalMemoryBytes,
			},
		},
	}
}

func (c *GRPCNodeClient) NextStatus() (*Status, error) {
	status, ok := <-c.statusQueue
	if !ok {
		return nil, io.EOF
	}
	return status, nil
}

func (c *GRPCNodeClient) startResultsProc() {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		ctx, cancl := context.WithTimeout(context.Background(), StreamConnectTimeout)
		defer cancl()
		stream, err := c.client.ProcessResults(ctx, &empty.Empty{})
		if err != nil {
			c.logger.Errorf("gRPCNodeClient: %v", err)
			return
		}
	mainLoop:
		for {
			result, err := stream.Recv()
			if err != io.EOF {
				close(c.resultsQueue)
				c.logger.Info("gRPCNodeClient: result stream ended.")
				return
			}
			if err != nil {
				c.logger.Errorf("gRPCNodeClient: %v", err)
				continue
			}
			select {
			case <-c.shutdown:
				break mainLoop
			default:
				c.resultsQueue <- &Result{
					CapturerUUID:  result.CapturerUuid,
					Recognized:    result.Recognized,
					TotalEntities: result.TotalEntities,
					RecognizedAt:  result.RecognizedAt.AsTime(),
					CapturedAt:    result.RecognizedAt.AsTime(),
				}
			}
		}
	}()
}

func (c *GRPCNodeClient) NextResult() (*Result, error) {
	result, ok := <-c.resultsQueue
	if !ok {
		return nil, io.EOF
	}
	return result, nil
}

func (c *GRPCNodeClient) Shutdown() error {
	close(c.shutdown)
	c.wg.Wait()
	if err := c.clientConn.Close(); err != nil {
		return fmt.Errorf("nodeclient: shutdown: %w", err)
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
