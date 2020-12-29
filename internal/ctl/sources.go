package ctl

import (
	"context"

	"github.com/eloylp/aton/internal/proto"
	"github.com/eloylp/aton/internal/video"
)

type Capturer interface {
	Start()
	NextOutput() (*video.Capture, error)
	Close()
	Status() string
	UUID() string
}

type DetectorClient interface {
	Connect() error
	LoadCategories(ctx context.Context, request *proto.LoadCategoriesRequest) (*proto.LoadCategoriesResponse, error)
	SendToRecognize(req *proto.RecognizeRequest) error
	NextRecognizeResponse() (*proto.RecognizeResponse, error)
	StartRecognize(ctx context.Context) error
	Shutdown() error
}
