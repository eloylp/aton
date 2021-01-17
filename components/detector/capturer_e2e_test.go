// +build e2e

package detector_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/eloylp/aton/pkg/test/helper"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/eloylp/aton/components/detector"
	"github.com/eloylp/aton/components/detector/metrics"
	"github.com/eloylp/aton/components/video"
)

func TestShutdownWhileBackPressured(t *testing.T) {
	// Prepare logger and metrics dependencies
	loggerOutput := bytes.NewBuffer(nil)
	logger := logrus.New()
	logger.SetOutput(loggerOutput)
	m := metrics.NewService("UUID")

	// Prepare the target handler
	sut := detector.NewCapturerHandler(logger, m, 10) // Make capturer handler the back pressured part.

	// Prepare video stream
	target := helper.ReplayedVideoStream(t, []string{faceBona1, faceBona2}, "/", 100)
	defer target.Close()

	capturer, err := video.NewMJPEGCapturer("UUID", target.URL, 50, logger)
	assert.NoError(t, err)
	sut.AddCapturer(capturer)
	time.Sleep(time.Second)
	sut.Shutdown()
}
