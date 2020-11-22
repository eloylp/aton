package grpc

import (
	"context"
	"io"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/eloylp/aton/internal/detector"
	"github.com/eloylp/aton/internal/detector/proto"
	"github.com/eloylp/aton/internal/logging"
)

type Service struct {
	detector detector.Classifier
	logging.Logger
	timeNow func() time.Time
}

func NewService(d detector.Classifier, logger logging.Logger, timeNow func() time.Time) *Service {
	return &Service{detector: d, Logger: logger, timeNow: timeNow}
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
			s.Logger.Info("ending detector consumer")
			return nil
		}
		if err != nil {
			s.Logger.Error(err)
			return err
		}
		var resp *proto.RecognizeResponse
		resp.CreatedAt = req.CreatedAt
		resp.Success = true
		cats, err := s.detector.FindCategories(req.Image)
		resp.RecognizedAt = timestamppb.Now()
		if err != nil {
			resp.Success = false
			resp.Message = err.Error()
			s.Logger.Error(err)
		}
		resp.Names = cats
		if err := server.Send(resp); err != nil {
			return err
		}
	}
}
