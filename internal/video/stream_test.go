// +build integration

package video_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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
	vc, err := video.NewMJPEGStreamCapturer(vs.URL, maxFrameBuffer)
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
	_, err := video.NewMJPEGStreamCapturer("tcp://127.0.0.1:8080", 5)
	assert.EqualError(t, err, "capturer (tcp://127.0.0.1:8080): only http or https scheme supported")
}

func TestOnCloseOutputChannelIsClosed(t *testing.T) {
	pictures := []string{faceBona1}
	vs := videoStream(t, pictures, "/")
	defer vs.Close()
	vc, err := video.NewMJPEGStreamCapturer(vs.URL, 1)
	assert.NoError(t, err)
	go vc.Start()
	time.Sleep(time.Second)
	vc.Close()
	_, closed := <-vc.Output()
	assert.True(t, closed)
}
