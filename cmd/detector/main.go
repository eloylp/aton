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
	logger.Infof("Starting detector service at %s", address)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		terminateAbnormally(logger, err)
	}
	s := grpc.NewServer()
	detector.RegisterDetectorServer(s, detector.NewServer(faceDetector, logger, time.Now))
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
		recvSig := <-ch
		logger.Infof("received shutdown signal %q, gracefully shutdown", recvSig.String())
		s.GracefulStop()
	}()
	if err := s.Serve(lis); err != nil {
		terminateAbnormally(logger, err)
	}
	logger.Infof("Stopped detector service at %s", address)
}

func terminateAbnormally(logger logging.Logger, err error) {
	logger.Error(err)
	os.Exit(1)
}
