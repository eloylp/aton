package detector

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

	"github.com/eloylp/aton/components/proto"
	"github.com/eloylp/aton/components/video"
)

type Service struct {
	UUID            string
	detector        Classifier
	capturerHandler *CapturerHandler
	logger          *logrus.Logger
	timeNow         func() time.Time
	L               *sync.Mutex
}

func NewService(
	uuid string, d Classifier,
	capturerHandler *CapturerHandler,
	logger *logrus.Logger,
	timeNow func() time.Time,
) *Service {
	return &Service{
		UUID:            uuid,
		detector:        d,
		capturerHandler: capturerHandler,
		logger:          logger,
		timeNow:         timeNow,
	}
}

func (s *Service) LoadCategories(_ context.Context, request *proto.LoadCategoriesRequest) (*empty.Empty, error) {
	if err := s.detector.SaveCategories(request.Categories, request.Image); err != nil {
		msg := fmt.Sprintf("LoadCategories(): error %v loading: %q", err, strings.Join(request.Categories, ","))
		s.logger.Error(msg)
		return nil, status.New(codes.Internal, msg).Err()
	}
	s.logger.Infof("LoadCategories(): loaded %q", strings.Join(request.Categories, ","))
	return nil, nil
}

func (s *Service) InformStatus(request *proto.InformStatusRequest, stream proto.Detector_InformStatusServer) error {
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

func (s *Service) ProcessResults(_ *empty.Empty, stream proto.Detector_ProcessResultsServer) error {
	for {
		capturerResult, err := s.capturerHandler.NextResult()
		if err == io.EOF {
			s.logger.Info("ProcessResults(): stopped, end of processing.")
			break
		}
		if err != nil {
			msg := fmt.Sprintf("ProcessResults(): %v", err)
			s.logger.Errorf(msg)
			return status.New(codes.Internal, msg).Err()
		}
		success := true
		cat, err := s.detector.FindCategories(capturerResult.Data)
		if err != nil {
			success = false
			msg := fmt.Sprintf("ProcessResults(): %v", err)
			s.logger.Errorf(msg)
		}
		recognizedAtProtoTime, err := ptypes.TimestampProto(s.timeNow())
		if err != nil {
			panic(err)
		}
		capturedAtProtoTime, err := ptypes.TimestampProto(capturerResult.Timestamp)
		if err != nil {
			panic(err)
		}
		err = stream.Send(&proto.Result{
			CapturerUuid: s.UUID,
			Matches:      cat,
			Success:      success,
			RecognizedAt: recognizedAtProtoTime,
			CapturedAt:   capturedAtProtoTime,
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
	err := s.capturerHandler.AddMJPEGCapturer(request.GetCapturerUuid(), request.GetCapturerUrl(), 10)
	if err != nil {
		return nil, status.New(codes.Internal, "addCapturer: "+err.Error()).Err()
	}
	return nil, nil
}

func (s *Service) RemoveCapturer(_ context.Context, request *proto.RemoveCapturerRequest) (*empty.Empty, error) {
	if err := s.capturerHandler.RemoveCapturer(request.GetCapturerUuid()); err != nil {
		msg := fmt.Sprintf("RemoveCapturer(): %v", err)
		s.logger.Error(msg)
		return nil, status.New(codes.NotFound, msg).Err() // TODO make type switch to gather not found err.
	}
	return nil, nil
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
	return &proto.Status{
		Description: "General status of detector",
		Capturers:   capt,
	}
}
