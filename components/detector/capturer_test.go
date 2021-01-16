// +build integration

package detector_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/eloylp/aton/components/detector"
	"github.com/eloylp/aton/components/detector/metrics"
	"github.com/eloylp/aton/components/video"
)

func TestProcessingTargetResults(t *testing.T) {
	// Prepare logger and metrics dependencies
	loggerOutput := bytes.NewBuffer(nil)
	logger := logrus.New()
	logger.SetOutput(loggerOutput)
	m := metrics.NewService("UUID")

	// Prepare the target handler
	sut := detector.NewCapturerHandler(logger, m, 100)

	// Prepare our test target, simulates a video capture.
	target := NewFakeTarget(t, []string{faceBona1, faceBona1})
	target.On("Start").Return()
	target.On("Close").Return()
	target.On("Status").Return(video.StatusRunning)
	target.On("UUID").Return("TEST")
	target.On("TargetURL").Return("http://example.com")
	// Including the target in our SUT, the target handler
	sut.AddCapturer(target)

	// Wait for the backbone to be filled with results
	assert.Eventually(t, func() bool {
		return sut.BackboneLen() == 2 // because we processed 2 images.
	}, time.Second, time.Millisecond)

	assert.Equal(t, []detector.CapturerStatus{
		{
			UUID:   "TEST",
			URL:    "http://example.com",
			Status: video.StatusRunning,
		},
	}, sut.Status())

	for i := 0; i < 2; i++ {
		data, err := sut.NextResult()
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)
		assert.Equal(t, 1301206, len(data.Data))
	}

	sut.Shutdown()

	target.AssertExpectations(t)

	// Assert logging
	logOutput := loggerOutput.String()
	assert.NotContains(t, logOutput, "level=error")
	assert.Contains(t, logOutput, "capturerHandler: added target with UUID: TEST")
	assert.Contains(t, logOutput, "capturerHandler: starting target with UUID: TEST")
	assert.Contains(t, logOutput, "capturerHandler: closing target with UUID: TEST")

	// Assert metrics
	rec := httptest.NewRecorder()
	m.HTTPHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))

	response := rec.Body.String()
	assert.Contains(t, response, `aton_detector_capturer_received_frames_total{uuid="TEST"} 2`)
	assert.Contains(t, response, `aton_detector_capturer_up{uuid="TEST"} 0`)
}

func TestRemoveCapturerFromHandler(t *testing.T) {
	// Prepare logger and metrics dependencies
	loggerOutput := bytes.NewBuffer(nil)
	logger := logrus.New()
	logger.SetOutput(loggerOutput)
	m := metrics.NewService("UUID")

	// Prepare the target handler
	sut := detector.NewCapturerHandler(logger, m, 100)

	// Prepare our test target, simulates a video capture.
	target := NewFakeTarget(t, []string{faceBona1, faceBona1})
	target.On("Start").Return()
	target.On("Close").Return()
	target.On("Status").Return(video.StatusRunning)
	target.On("UUID").Return("UUID")
	target.On("TargetURL").Return("http://example.com")
	// Including the target in our SUT, the target handler
	sut.AddCapturer(target)

	// We remove the capturer and check count in status.
	capt, err := sut.RemoveCapturer("UUID")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(sut.Status()))
	assert.Equal(t, "UUID", capt.UUID())
	assert.Equal(t, "http://example.com", capt.TargetURL())

	// Check that now throws error because does not exist.
	_, err = sut.RemoveCapturer("UUID")
	assert.EqualError(t, err, "capturerHandler: capturer with UUID UUID not found")
}
