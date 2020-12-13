// +build integration

package ctl_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/eloylp/aton/internal/proto"
	"github.com/eloylp/aton/internal/video"
	"github.com/eloylp/aton/pkg/test/helper"
)

type fakeDetectorClient struct {
	mock.Mock
	recognizeReq     chan *proto.RecognizeRequest
	recognizeResp    chan *proto.RecognizeResponse
	errPercent       float32
	processedSuccess float32
	processedFailed  float32
}

func (f *fakeDetectorClient) SendToRecognize(req *proto.RecognizeRequest) error {
	f.recognizeReq <- req
	return f.Called(req).Error(0)
}

func (f *fakeDetectorClient) NextRecognizeResponse() (*proto.RecognizeResponse, error) {
	resp, ok := <-f.recognizeResp
	if !ok {
		return nil, io.EOF
	}
	return resp, nil
}

func (f *fakeDetectorClient) StartRecognize(ctx context.Context) error {
	args := f.Called(ctx)
	go func() {
		for req := range f.recognizeReq {
			if f.processedFailed/(f.processedSuccess+f.processedFailed) < f.errPercent/100 {
				f.recognizeResp <- &proto.RecognizeResponse{
					Names:     []string{},
					Success:   false,
					Message:   "Not recognized",
					CreatedAt: req.CreatedAt,
				}
				f.processedFailed++
				continue
			}
			f.recognizeResp <- &proto.RecognizeResponse{
				Names:     []string{"bona"},
				Success:   true,
				Message:   "Matches on pictures",
				CreatedAt: req.CreatedAt,
			}
			f.processedSuccess++
		}
	}()
	return args.Error(0)
}

func newFakeDetectorClient(errPercent float32) *fakeDetectorClient {
	return &fakeDetectorClient{
		recognizeReq:  make(chan *proto.RecognizeRequest, 10),
		recognizeResp: make(chan *proto.RecognizeResponse, 10),
		errPercent:    errPercent,
	}
}

func (f *fakeDetectorClient) Connect() error {
	args := f.Called()
	return args.Error(0)
}

func (f *fakeDetectorClient) LoadCategories(
	ctx context.Context,
	request *proto.LoadCategoriesRequest) (*proto.LoadCategoriesResponse, error) {
	args := f.Called(ctx, request)
	return args.Get(0).(*proto.LoadCategoriesResponse), args.Error(1)
}

func (f *fakeDetectorClient) Shutdown() error {
	close(f.recognizeReq)
	close(f.recognizeResp)
	return f.Called().Error(0)
}

type fakeCapturer struct {
	mock.Mock
	stream    chan *video.Capture
	uuidIdent string
}

func (f *fakeCapturer) UUID() string {
	return f.uuidIdent
}

func newFakeCapturer(t *testing.T, uuid string, images []string) *fakeCapturer {
	stream := make(chan *video.Capture, 10)
	for _, img := range images {
		stream <- &video.Capture{
			Data:      helper.ReadFile(t, img),
			Timestamp: time.Now(),
		}
	}
	return &fakeCapturer{
		uuidIdent: uuid,
		stream:    stream,
	}
}

func (f *fakeCapturer) Start() {
	f.Called()
}

func (f *fakeCapturer) NextOutput() (*video.Capture, error) {
	capt, ok := <-f.stream
	if !ok {
		return nil, io.EOF
	}
	return capt, nil
}

func (f *fakeCapturer) Close() {
	close(f.stream)
	f.Called()
}

func (f *fakeCapturer) Status() string {
	args := f.Called()
	return args.String(0)
}
