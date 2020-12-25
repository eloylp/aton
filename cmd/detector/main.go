package main

import (
	"os"
	"time"

	"github.com/eloylp/aton/internal/detector"
	"github.com/eloylp/aton/internal/detector/grpc"
	"github.com/eloylp/aton/internal/logging"
)

const (
	DetectorListenAddress  = "DETECTOR_ADDR"
	DetectorMetricsAddress = "DETECTOR_METRICS_ADDR"
	DetectorModelDir       = "DETECTOR_MODEL_DIR"
)

func main() {
	address := os.Getenv(DetectorListenAddress)
	metricsAddress := os.Getenv(DetectorMetricsAddress)
	modelDir := os.Getenv(DetectorModelDir)
	logger := logging.NewBasicLogger(os.Stdout)
	faceDetector, err := detector.NewGoFaceDetector(modelDir)
	if err != nil {
		terminateAbnormally(logger, err)
	}
	service := grpc.NewService(faceDetector, logger, time.Now)
	server := grpc.NewServer(address, service, metricsAddress, logger)
	if err := server.Start(); err != nil {
		terminateAbnormally(logger, err)
	}
}

func terminateAbnormally(logger logging.Logger, err error) {
	logger.Error(err)
	os.Exit(1)
}
