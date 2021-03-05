package node

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/eloylp/aton/components/node/metrics"
	"github.com/eloylp/aton/components/node/system"
	"github.com/eloylp/aton/components/proto"
	"github.com/eloylp/aton/components/video"
)

type Service struct {
	node            Classifier
	capturerHandler *CapturerHandler
	metricsService  *metrics.Service
	logger          *logrus.Logger
	timeNow         func() time.Time
	L               *sync.Mutex
}

func NewService(
	d Classifier,
	capturerHandler *CapturerHandler,
	metricsService *metrics.Service,
	logger *logrus.Logger,
	timeNow func() time.Time,
) *Service {
	return &Service{
		node:            d,
		capturerHandler: capturerHandler,
		metricsService:  metricsService,
		logger:          logger,
		timeNow:         timeNow,
	}
}

func (s *Service) LoadCategories(_ context.Context, request *proto.LoadCategoriesRequest) (*empty.Empty, error) {
	if err := s.node.SaveCategories(request.Categories, request.Image); err != nil {
		msg := fmt.Sprintf("LoadCategories(): error %v loading: %q", err, strings.Join(request.Categories, ","))
		s.logger.Error(msg)
		return nil, status.New(codes.Internal, msg).Err()
	}
	s.logger.Infof("LoadCategories(): loaded %q", strings.Join(request.Categories, ","))
	return &empty.Empty{}, nil
}

func (s *Service) InformStatus(request *proto.InformStatusRequest, stream proto.Node_InformStatusServer) error {
	for {
		err := stream.Send(s.Status())
		if err == io.EOF {
			s.logger.Info("InformStatus(): stopped by client.")
			break
		}
		if err != nil {
			s.logger.Errorf("InformStatus(): %v", err.Error())
			return err
		}
		time.Sleep(request.Interval.AsDuration())
	}
	return nil
}

func (s *Service) ProcessResults(_ *empty.Empty, stream proto.Node_ProcessResultsServer) error {
	for {
		capturerResult, err := s.capturerHandler.NextResult()
		if err == io.EOF {
			s.logger.Info("ProcessResults(): stopped, end of processing.")
			break
		}
		s.metricsService.IncProcessedFramesTotal()
		if err != nil {
			msg := fmt.Sprintf("ProcessResults(): %v", err)
			s.logger.Error(msg)
			s.metricsService.IncFailedFramesTotal()
			return status.New(codes.Internal, msg).Err()
		}
		now := s.timeNow()
		resp, err := s.node.FindCategories(capturerResult.Data)
		if err != nil {
			msg := fmt.Sprintf("ProcessResults(): %v", err)
			s.logger.Error(msg)
			s.metricsService.IncFailedFramesTotal()
			return status.New(codes.Internal, msg).Err()
		}
		elapsed := time.Since(now)
		s.metricsService.IncProcessedTimeFramesSeconds(elapsed.Seconds())
		s.metricsService.AddEntitiesTotal(resp.TotalEntities)
		s.metricsService.AddUnrecognizedEntitiesTotal(resp.TotalEntities - len(resp.Matches))
		recognizedAtProtoTime, err := ptypes.TimestampProto(s.timeNow())
		if err != nil {
			panic(err)
		}
		capturedAtProtoTime, err := ptypes.TimestampProto(capturerResult.Timestamp)
		if err != nil {
			panic(err)
		}
		err = stream.Send(&proto.Result{
			Recognized:    resp.Matches,
			TotalEntities: int32(resp.TotalEntities),
			RecognizedAt:  recognizedAtProtoTime,
			CapturedAt:    capturedAtProtoTime,
		})
		if err == io.EOF {
			s.logger.Info("ProcessResults(): stopped, end of processing.")
			break
		}
		if err != nil {
			msg := fmt.Sprintf("ProcessResults(): %v", err)
			s.logger.Errorf(msg)
			return status.New(codes.Internal, msg).Err()
		}
	}
	return nil
}

func (s *Service) AddCapturer(_ context.Context, request *proto.AddCapturerRequest) (*empty.Empty, error) {
	err := s.capturerHandler.AddMJPEGCapturer(request.GetCapturerUuid(), request.GetCapturerUrl(), 10) // Todo ...
	if err != nil {
		return nil, status.New(codes.Internal, "addCapturer: "+err.Error()).Err()
	}
	s.metricsService.CapturerUP(request.CapturerUuid, request.CapturerUrl)
	return &empty.Empty{}, nil
}

func (s *Service) RemoveCapturer(_ context.Context, request *proto.RemoveCapturerRequest) (*empty.Empty, error) {
	capt, err := s.capturerHandler.RemoveCapturer(request.GetCapturerUuid())
	if err != nil {
		msg := fmt.Sprintf("RemoveCapturer(): %v", err)
		s.logger.Error(msg)
		return nil, status.New(codes.NotFound, msg).Err() // TODO make type switch to gather not found err.
	}
	s.metricsService.CapturerDown(capt.UUID(), capt.TargetURL())
	return &empty.Empty{}, nil
}

func (s *Service) Status() *proto.Status {
	cs := s.capturerHandler.Status()
	capt := make([]*proto.Capturer, len(cs))

	for _, c := range s.capturerHandler.Status() {
		var pStatus proto.CapturerStatus
		if c.Status == video.StatusRunning {
			pStatus = proto.CapturerStatus_CAPTURER_STATUS_OK
		}
		if c.Status == video.StatusNotRunning {
			pStatus = proto.CapturerStatus_CAPTURER_STATUS_CONNECTION_RETRY
		}
		capt = append(capt, &proto.Capturer{
			Uuid:   c.UUID,
			Url:    c.URL,
			Status: pStatus,
		})
	}
	memory := system.Memory()
	load := system.LoadAverage()
	network := system.Network()

	return &proto.Status{
		Description: "General status of node",
		Capturers:   capt,
		System: &proto.System{
			CpuCount: int32(system.CPUCount()),
			Memory: &proto.Memory{
				TotalMemoryBytes: memory.TotalBytes,
				UsedMemoryBytes:  memory.UsedBytes,
			},
			LoadAverage: &proto.LoadAverage{
				Avg_1:  load.LoadAvg1,
				Avg_5:  load.LoadAvg5,
				Avg_15: load.LoadAvg15,
			},
			Network: &proto.Network{
				TxBytesSec: network.TxBytesSec,
				RxBytesSec: network.RxBytesSec,
			},
		},
	}
}

func (s *Service) Shutdown() {
	s.capturerHandler.Shutdown()
}
