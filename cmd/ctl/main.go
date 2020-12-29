package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/eloylp/aton/internal/ctl/config"
	"github.com/eloylp/aton/internal/ctl/engine"
	"github.com/eloylp/aton/internal/ctl/grpc"
	"github.com/eloylp/aton/internal/ctl/metrics"
	"github.com/eloylp/aton/internal/video"
)

func main() {
	logger := logrus.New()
	metricsService := metrics.NewService()
	dc := grpc.NewDetectorClient("127.0.0.1:8082", logger, metricsService)
	c := engine.New(dc, metricsService, logger,
		config.WithListenAddress("0.0.0.0:8081"),
	)
	capturer, err := video.NewMJPEGCapturer("capt1", os.Getenv("CAPT_URL"), 10, logger)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	c.AddCapturer(capturer)
	if err := c.Start(); err != nil {
		logger.Error(err)
	}
}
