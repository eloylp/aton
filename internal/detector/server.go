package detector

import (
	"context"
	"io"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/eloylp/aton/internal/logging"
)

type Server struct {
	detector Facial
	logging.Logger
	timeNow func() time.Time
}

func NewServer(detector Facial, logger logging.Logger, timeNow func() time.Time) *Server {
	return &Server{detector: detector, Logger: logger, timeNow: timeNow}
}

func (s *Server) LoadCategories(_ context.Context, r *LoadCategoriesRequest) (*LoadCategoriesResponse, error) {
	if err := s.detector.SaveFaces(r.Categories, r.Image); err != nil {
		return nil, err
	}
	return &LoadCategoriesResponse{
		Success: true,
		Message: "categories loaded",
	}, nil
}

func (s *Server) Recognize(server Detector_RecognizeServer) error {
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
