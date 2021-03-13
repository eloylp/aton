// +build integration

package video_test

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/eloylp/aton/components/video"
	"github.com/eloylp/aton/pkg/test/helper"
)

// TODO. this is the beginning of the implementation of the video
// TODO. capturer. The below tests must be decoupled from implementation
// TODO. via interface, but since we still dont know if the current
// TODO. interface is the final one, we are delaying this refactor.

var (
	imagesDir = "../../samples/images"
	faceBona1 = filepath.Join(imagesDir, "bona.jpg")
	faceBona2 = filepath.Join(imagesDir, "bona2.jpg")
	faceBona3 = filepath.Join(imagesDir, "bona3.jpg")
	faceBona4 = filepath.Join(imagesDir, "bona4.jpg")
)

func TestCapture(t *testing.T) {
	pictures := []string{faceBona1, faceBona2, faceBona3, faceBona4}
	vs := helper.VideoStream(t, pictures, "/")
	defer vs.Close()
	maxFrameBuffer := 100
	logger := testLogger(bytes.NewBuffer(nil))
	vc, err := video.NewMJPEGCapturer("uuid", vs.URL, maxFrameBuffer, logger)
	assert.NoError(t, err)
	go vc.Start()
	time.AfterFunc(500*time.Millisecond, func() {
		vc.Close()
	})
	var i int
	for {
		k, err := vc.NextOutput()
		if err != nil {
			break
		}
		expected := helper.ReadFile(t, pictures[i])
		got := k.Data
		assert.Equal(t, "uuid", k.CapturerUUID)
		assert.Equal(t, expected, got, "Image no %v does not match", i)
		i++
	}
}

func testLogger(output io.Writer) *logrus.Logger {
	l := logrus.New()
	l.SetOutput(output)
	return l
}

func TestNonSupportedURLScheme(t *testing.T) {
	logger := testLogger(bytes.NewBuffer(nil))
	_, err := video.NewMJPEGCapturer("uuid", "tcp://127.0.0.1:8080", 5, logger)
	assert.EqualError(t, err, "capturer (tcp://127.0.0.1:8080): only http or https scheme supported")
}

func TestOnCloseOutputChannelIsClosed(t *testing.T) {
	pictures := []string{faceBona1}
	vs := helper.VideoStream(t, pictures, "/")
	defer vs.Close()
	logger := testLogger(bytes.NewBuffer(nil))
	vc, err := video.NewMJPEGCapturer("uuid", vs.URL, 1, logger)
	assert.NoError(t, err)
	go vc.Start()
	time.Sleep(time.Second)
	vc.Close()
	_, _ = vc.NextOutput()   // exhaust
	_, err = vc.NextOutput() // exhaust
	assert.EqualError(t, err, io.EOF.Error())
}

func TestErrorConnectionRefusedLogged(t *testing.T) {
	w := bytes.NewBuffer(nil)
	logger := testLogger(w)
	vc, err := video.NewMJPEGCapturer("uuid", "http://127.0.0.2", 1, logger)
	assert.NoError(t, err)
	go vc.Start()
	defer vc.Close()
	time.Sleep(100 * time.Millisecond)
	assert.NoError(t, err)
	assert.Contains(t, w.String(), "capturer: ")
	assert.Contains(t, w.String(), "connect: connection refused")
}

func TestCloseWorksEvenDuringProcessingFrames(t *testing.T) {
	framesInFlight := 1000
	pictures := make([]string, framesInFlight)
	for i := 0; i < framesInFlight; i++ {
		pictures[i] = faceBona1
	}
	vs := helper.VideoStream(t, pictures, "/")
	defer vs.Close()
	logger := testLogger(bytes.NewBuffer(nil))
	vc, err := video.NewMJPEGCapturer("uuid", vs.URL, framesInFlight, logger)
	assert.NoError(t, err)
	go vc.Start()
	time.Sleep(500 * time.Millisecond)
	vc.Close()
	time.Sleep(50 * time.Millisecond)
	_, err = vc.NextOutput()
	assert.NoError(t, err, "output channel needs to have residual frames")
	count := 0
	for {
		_, err := vc.NextOutput()
		if err == io.EOF {
			break
		}
		count++
	}
	t.Logf("Output channel residual length: %v", count)
}

func TestInitialState(t *testing.T) {
	w := bytes.NewBuffer(nil)
	logger := testLogger(w)
	vc, err := video.NewMJPEGCapturer("uuid", "http://localhost:9999", 10, logger)
	assert.NoError(t, err)
	assert.Equal(t, video.StatusNotRunning, vc.Status(), "Init state must be %s", video.StatusNotRunning)
}

func TestRunningState(t *testing.T) {
	pictures := []string{faceBona1, faceBona2, faceBona3, faceBona4}
	vs := helper.VideoStream(t, pictures, "/")
	defer vs.Close()
	w := bytes.NewBuffer(nil)
	logger := testLogger(w)
	vc, err := video.NewMJPEGCapturer("uuid", vs.URL, 10, logger)
	assert.NoError(t, err)
	go vc.Start()
	defer vc.Close()
	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, video.StatusRunning, vc.Status(), "Running state must be %s", video.StatusRunning)
}

func TestExpBackoffReconnectPeriods(t *testing.T) {
	t.Parallel()
	w := bytes.NewBuffer(nil)
	logger := testLogger(w)
	expectedConnections := 5
	netListener, netServer := helper.NetCloserServer(t, expectedConnections)
	defer netListener.Close()
	vc, err := video.NewMJPEGCapturer("uuid", "http://"+netListener.Addr().String(), 10, logger)
	assert.NoError(t, err)
	startTime := time.Now()
	go vc.Start()
	time.Sleep(40 * time.Second)
	vc.Close()
	resultConnections := make([]time.Time, expectedConnections)
	assert.Len(t, netServer, expectedConnections)
	var i int
	for connTime := range netServer {
		resultConnections[i] = connTime
		i++
	}
	marginErr := 50 * time.Millisecond
	assert.WithinDuration(t, startTime, resultConnections[0], marginErr)
	assert.WithinDuration(t, startTime.Add(2*time.Second), resultConnections[1], marginErr)
	assert.WithinDuration(t, startTime.Add(6*time.Second), resultConnections[2], marginErr)
	assert.WithinDuration(t, startTime.Add(14*time.Second), resultConnections[3], marginErr)
	assert.WithinDuration(t, startTime.Add(30*time.Second), resultConnections[4], marginErr)
}

func TestShutdownWhileBackPressured(t *testing.T) {
	pictures := []string{faceBona1, faceBona2, faceBona3, faceBona4}
	vs := helper.ReplayedVideoStream(t, pictures, "/", 100, 25)
	defer vs.Close()
	w := bytes.NewBuffer(nil)
	logger := testLogger(w)
	vc, err := video.NewMJPEGCapturer("uuid", vs.URL, 10, logger)
	assert.NoError(t, err)
	time.AfterFunc(time.Second, func() {
		vc.Close()
	})
	vc.Start()
	fmt.Println(w.String())
}
