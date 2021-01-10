// +build integration

package detector_test

import (
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/eloylp/aton/components/video"
	"github.com/eloylp/aton/pkg/test/helper"
)

type FakeTarget struct {
	mock.Mock
	output chan *video.Capture
}

func NewFakeTarget(t *testing.T, outputs []string) *FakeTarget {
	output := make(chan *video.Capture, len(outputs))
	for _, o := range outputs {
		output <- &video.Capture{
			Data:      helper.ReadFile(t, o),
			Timestamp: time.Now(),
		}
	}
	return &FakeTarget{output: output}
}

func (ft *FakeTarget) Start() {
	ft.Called()
}

func (ft *FakeTarget) NextOutput() (*video.Capture, error) {
	c, ok := <-ft.output
	if !ok {
		return nil, io.EOF
	}
	return c, nil
}

func (ft *FakeTarget) Close() {
	close(ft.output)
	ft.Called()
}

func (ft *FakeTarget) Status() string {
	return ft.Called().String(0)
}

func (ft *FakeTarget) UUID() string {
	return ft.Called().String(0)
}

func (ft *FakeTarget) TargetURL() string {
	return ft.Called().String(0)
}
