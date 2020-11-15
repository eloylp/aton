package detector

import (
	"context"
	"io"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/eloylp/aton/internal/logging"
)

type GRPCService struct {
	detector Facial
	logging.Logger
	timeNow func() time.Time
}

func NewGRPCService(detector Facial, logger logging.Logger, timeNow func() time.Time) *GRPCService {
	return &GRPCService{detector: detector, Logger: logger, timeNow: timeNow}
}

func (s *GRPCService) LoadCategories(_ context.Context, r *LoadCategoriesRequest) (*LoadCategoriesResponse, error) {
	if err := s.detector.SaveFaces(r.Categories, r.Image); err != nil {
		return nil, err
	}
	return &LoadCategoriesResponse{
		Success: true,
		Message: "categories loaded",
	}, nil
}

func (s *GRPCService) Recognize(server Detector_RecognizeServer) error {
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
		var resp *RecognizeResponse
		resp.CreatedAt = req.CreatedAt
		resp.Success = true
		cats, err := s.detector.FindFaces(req.Image)
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
