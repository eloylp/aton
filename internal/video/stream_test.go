// +build integration

package video_test

import (
	"bytes"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/eloylp/aton/internal/logging"
	"github.com/eloylp/aton/internal/video"
)

var (
	imagesDir = "../../images"
	faceBona1 = filepath.Join(imagesDir, "bona.jpg")
	faceBona2 = filepath.Join(imagesDir, "bona2.jpg")
	faceBona3 = filepath.Join(imagesDir, "bona3.jpg")
	faceBona4 = filepath.Join(imagesDir, "bona4.jpg")
)

func TestCapture(t *testing.T) {
	pictures := []string{faceBona1, faceBona2, faceBona3, faceBona4}
	vs := videoStream(t, pictures, "/")
	defer vs.Close()
	maxFrameBuffer := 100
	logger := logging.NewBasicLogger(bytes.NewBuffer(nil))
	vc, err := video.NewMJPEGStreamCapturer(vs.URL, maxFrameBuffer, logger)
	assert.NoError(t, err)
	go vc.Start()
	output := vc.Output()
	time.AfterFunc(500*time.Millisecond, func() {
		vc.Close()
	})
	var i int
	for k := range output {
		expected := readFile(t, pictures[i])
		got := k.Data
		assert.Equal(t, expected, got, "Image no %v does not match", i)
		i++
	}
}

func TestNonSupportedURLScheme(t *testing.T) {
	logger := logging.NewBasicLogger(bytes.NewBuffer(nil))
	_, err := video.NewMJPEGStreamCapturer("tcp://127.0.0.1:8080", 5, logger)
	assert.EqualError(t, err, "capturer (tcp://127.0.0.1:8080): only http or https scheme supported")
}

func TestOnCloseOutputChannelIsClosed(t *testing.T) {
	pictures := []string{faceBona1}
	vs := videoStream(t, pictures, "/")
	defer vs.Close()
	logger := logging.NewBasicLogger(bytes.NewBuffer(nil))
	vc, err := video.NewMJPEGStreamCapturer(vs.URL, 1, logger)
	assert.NoError(t, err)
	go vc.Start()
	time.Sleep(time.Second)
	vc.Close()
	_, closed := <-vc.Output()
	assert.True(t, closed)
}

func TestErrorConnectionRefusedLogged(t *testing.T) {
	w := bytes.NewBuffer(nil)
	logger := logging.NewBasicLogger(w)
	vc, err := video.NewMJPEGStreamCapturer("http://127.0.0.2", 1, logger)
	assert.NoError(t, err)
	go vc.Start()
	time.Sleep(100 * time.Millisecond)
	assert.NoError(t, err)
	assert.Contains(t, w.String(), "capturer: ")
	assert.Contains(t, w.String(), "connect: connection refused")
}
