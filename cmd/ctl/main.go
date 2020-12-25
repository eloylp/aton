package main

import (
	"os"

	"github.com/eloylp/aton/internal/ctl"
	"github.com/eloylp/aton/internal/ctl/config"
	"github.com/eloylp/aton/internal/ctl/grpc"
	"github.com/eloylp/aton/internal/logging"
	"github.com/eloylp/aton/internal/video"
)

func main() {
	logger := logging.NewBasicLogger(os.Stdout)
	dc := grpc.NewDetectorClient("127.0.0.1:8082", logger)
	c := ctl.New(dc,
		config.WithListenAddress("0.0.0.0:8081"),
		config.WithLoggerOutput(os.Stdout),
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
