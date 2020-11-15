package main

import (
	"os"
	"time"

	"github.com/eloylp/aton/internal/detector"
	"github.com/eloylp/aton/internal/detector/grpc"
	"github.com/eloylp/aton/internal/logging"
)

const (
	DetectorListenAddress = "DETECTOR_ADDR"
	DetectorModelDir      = "DETECTOR_MODEL_DIR"
)

func main() {
	address := os.Getenv(DetectorListenAddress)
	modelDir := os.Getenv(DetectorModelDir)
	logger := logging.NewBasicLogger(os.Stdout)
	faceDetector, err := detector.NewGoFaceDetector(modelDir)
	if err != nil {
		terminateAbnormally(logger, err)
	}
	service := grpc.NewService(faceDetector, logger, time.Now)
	server := grpc.NewServer(service, logger, address)
	if err := server.Start(); err != nil {
		terminateAbnormally(logger, err)
	}
}

func terminateAbnormally(logger logging.Logger, err error) {
	logger.Error(err)
	os.Exit(1)
}
