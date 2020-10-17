// +build integration

package video_test

import (
	"github.com/eloylp/aton/internal/video"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
	"time"
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
	vc := video.NewMJPEGStreamCapturer(vs.URL, maxFrameBuffer)
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
