package main

import (
	"os"
	"time"

	"github.com/eloylp/aton/internal/detector"
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
	service := detector.NewGRPCService(faceDetector, logger, time.Now)
	server := detector.NewGRPCServer(service, logger, address)
	if err := server.Start(); err != nil {
		terminateAbnormally(logger, err)
	}
}

func terminateAbnormally(logger logging.Logger, err error) {
	logger.Error(err)
	os.Exit(1)
}
