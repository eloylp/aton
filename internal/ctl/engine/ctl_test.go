// +build integration

package engine_test

import (
	"bytes"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/eloylp/aton/internal/ctl/config"
	"github.com/eloylp/aton/internal/ctl/engine"
	"github.com/eloylp/aton/internal/ctl/metrics"
)

var (
	imagesDir = "../../../samples/images"
	faceBona1 = filepath.Join(imagesDir, "bona.jpg")
	faceBona2 = filepath.Join(imagesDir, "bona2.jpg")
	faceBona3 = filepath.Join(imagesDir, "bona3.jpg")
	faceBona4 = filepath.Join(imagesDir, "bona4.jpg")
	allFaces  = []string{faceBona1, faceBona2, faceBona3, faceBona4}
)

const (
	ctlListenAddress = "0.0.0.0:8080"
)

func TestCtlDoesBasicFlow(t *testing.T) {
	loggerOutput := bytes.NewBuffer(nil)
	logger := testLogger(loggerOutput)
	dc := newFakeDetectorClient(25)
	dc.On("Connect").Return(nil)
	dc.On("StartRecognize", mock.Anything).Return(nil)
	dc.On("SendToRecognize", mock.Anything).Return(nil)
	dc.On("Shutdown").Return(nil)
	sutCTL := engine.New(
		dc,
		metrics.NewService(),
		logger,
		config.WithListenAddress(ctlListenAddress),
		config.WithDetectors(config.Detector{
			Address: "127.0.0.1:8080",
			UUID:    "09AF",
		}),
	)

	go func() {
		assert.NoError(t, sutCTL.Start())
	}()

	capturer := newFakeCapturer(t, "cap1", allFaces)
	capturer.On("Start").Return()
	capturer.On("Close").Return()
	sutCTL.AddCapturer(capturer)
	time.Sleep(500 * time.Millisecond)

	metricsData := string(fetchResource(t, "http://"+ctlListenAddress+"/metrics"))

	assert.Contains(t, metricsData, `aton_ctl_detector_up{uuid="09AF"} 1`)
	assert.Contains(t, metricsData, `aton_ctl_capturer_up{uuid="cap1"} 1`)

	assert.Contains(t, metricsData, `aton_ctl_capturer_received_frames_total{uuid="cap1"} 4`)
	assert.NotContains(t, metricsData, `aton_ctl_capturer_failed_frames_total{uuid="cap1"}`)
	assert.Contains(t, metricsData, `aton_ctl_processed_frames_total{uuid="09AF"} 4`)
	assert.Contains(t, metricsData, `aton_ctl_unrecognized_frames_total{uuid="09AF"} 1`)

	assert.Contains(t, loggerOutput.String(), "initializeResultProcessor(): detected: bona")
	assert.Contains(t, loggerOutput.String(), "not detected:")
	assert.NotContains(t, loggerOutput.String(), "level=error")

	sutCTL.Shutdown()
	dc.AssertExpectations(t)
}
