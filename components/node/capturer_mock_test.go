// +build integration

package node_test

import (
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/eloylp/aton/components/video"
	"github.com/eloylp/aton/pkg/test/helper"
)

type FakeCapturer struct {
	mock.Mock
	output chan *video.Capture
}

func NewFakeTarget(t *testing.T, outputs []string) *FakeCapturer {
	output := make(chan *video.Capture, len(outputs))
	for _, o := range outputs {
		output <- &video.Capture{
			Data:      helper.ReadFile(t, o),
			Timestamp: time.Now(),
		}
	}
	return &FakeCapturer{output: output}
}

func (ft *FakeCapturer) Start() {
	ft.Called()
}

func (ft *FakeCapturer) NextOutput() (*video.Capture, error) {
	c, ok := <-ft.output
	if !ok {
		return nil, io.EOF
	}
	return c, nil
}

func (ft *FakeCapturer) Close() {
	close(ft.output)
	ft.Called()
}

func (ft *FakeCapturer) Status() string {
	return ft.Called().String(0)
}

func (ft *FakeCapturer) UUID() string {
	return ft.Called().String(0)
}

func (ft *FakeCapturer) TargetURL() string {
	return ft.Called().String(0)
}
