package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/eloylp/aton/internal/detector/engine"
	"github.com/eloylp/aton/internal/detector/grpc"
	"github.com/eloylp/aton/internal/logging"
)

const (
	DetectorUUID           = "DETECTOR_UUID"
	DetectorListenAddress  = "DETECTOR_ADDR"
	DetectorMetricsAddress = "DETECTOR_METRICS_ADDR"
	DetectorModelDir       = "DETECTOR_MODEL_DIR"
)

func main() {
	UUID := os.Getenv(DetectorUUID)
	address := os.Getenv(DetectorListenAddress)
	metricsAddress := os.Getenv(DetectorMetricsAddress)
	modelDir := os.Getenv(DetectorModelDir)
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	faceDetector, err := engine.NewGoFaceDetector(modelDir)
	if err != nil {
		terminateAbnormally(logger, err)
	}
	service := grpc.NewService(UUID, faceDetector, logger, time.Now)
	server := grpc.NewServer(address, service, metricsAddress, logger)
	if err := server.Start(); err != nil {
		terminateAbnormally(logger, err)
	}
}

func terminateAbnormally(logger logging.Logger, err error) {
	logger.Error(err)
	os.Exit(1)
}
