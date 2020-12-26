package grpc

import (
	"context"
	"io"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/eloylp/aton/internal/detector"
	"github.com/eloylp/aton/internal/proto"
)

type Service struct {
	UUID     string
	detector detector.Classifier
	logger   *logrus.Logger
	timeNow  func() time.Time
}

func NewService(uuid string, d detector.Classifier, logger *logrus.Logger, timeNow func() time.Time) *Service {
	return &Service{
		UUID:     uuid,
		detector: d,
		logger:   logger,
		timeNow:  timeNow,
	}
}

func (s *Service) LoadCategories(_ context.Context, r *proto.LoadCategoriesRequest) (*proto.LoadCategoriesResponse, error) {
	if err := s.detector.SaveCategories(r.Categories, r.Image); err != nil {
		return nil, err
	}
	return &proto.LoadCategoriesResponse{
		Success: true,
		Message: "categories loaded",
	}, nil
}

func (s *Service) Recognize(server proto.Detector_RecognizeServer) error {
	for {
		req, err := server.Recv()
		if err == io.EOF {
			s.logger.Info("ending detector consumer")
			return nil
		}
		if err != nil {
			s.logger.Error(err)
			return err
		}
		resp := &proto.RecognizeResponse{}
		resp.CreatedAt = req.CreatedAt
		resp.ProcessedBy = s.UUID
		resp.Success = true
		cats, err := s.detector.FindCategories(req.Image)
		resp.RecognizedAt = timestamppb.Now()
		if err != nil {
			resp.Success = false
			resp.Message = err.Error()
			s.logger.Error(err)
		}
		resp.Names = cats
		if err := server.Send(resp); err != nil {
			return err
		}
	}
}
