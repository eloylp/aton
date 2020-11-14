package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/eloylp/aton/internal/detector"
	"github.com/eloylp/aton/internal/logging"
)

func main() {
	address := os.Getenv("DETECTOR_ADDR")
	modelDir := os.Getenv("DETECTOR_MODELS_DIR")
	logger := logging.NewBasicLogger(os.Stdout)
	faceDetector, err := detector.NewGoFaceDetector(modelDir)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	logger.Infof("Starting detector service at %s", address)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	s := grpc.NewServer()
	detector.RegisterDetectorServer(s, detector.NewServer(faceDetector, logger, time.Now))
	if err := s.Serve(lis); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	logger.Infof("Stopped detector service at %s", address)
}
