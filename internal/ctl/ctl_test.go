// +build integration

package ctl_test

import (
	"bytes"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/eloylp/aton/internal/ctl"
	"github.com/eloylp/aton/internal/ctl/config"
)

var (
	imagesDir = "../../samples/images"
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
	var loggerOutput bytes.Buffer
	dc := newFakeDetectorClient(25)
	dc.On("Connect").Return(nil)
	dc.On("StartRecognize", mock.Anything).Return(nil)
	dc.On("SendToRecognize", mock.Anything).Return(nil)
	dc.On("Shutdown").Return(nil)
	sutCTL := ctl.New(
		dc,
		config.WithListenAddress(ctlListenAddress),
		config.WithLoggerOutput(&loggerOutput),
		config.WithDetectors(config.Detector{
			Address: "127.0.0.1",
			UUID:    "09AF",
		}),
	)

	go func() {
		assert.NoError(t, sutCTL.Start())
	}()

	capturer := newFakeCapturer(t, "cap1", allFaces)
	capturer.On("Start").Return()
	capturer.On("Close").Return()
	err := sutCTL.AddCapturer(capturer)
	assert.NoError(t, err)
	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, int32(1), sutCTL.Stats().CurrentDetectors(), "Unexpected detectors number")
	assert.Equal(t, int32(1), sutCTL.Stats().CurrentCapturers(), "Unexpected capturers number")

	assert.Equal(t, int64(4), sutCTL.Stats().Processed(), "Unexpected total processed frames number")
	assert.Equal(t, int64(3), sutCTL.Stats().ProcessedSuccess(), "Unexpected success processed frames number")
	assert.Equal(t, int64(1), sutCTL.Stats().ProcessedFailed(), "Unexpected failed processed frames number")

	assert.Contains(t, loggerOutput.String(), "detected: bona")
	assert.NotContains(t, loggerOutput.String(), "level=error")

	sutCTL.Shutdown()
	assert.NoError(t, err)
	dc.AssertExpectations(t)
}
