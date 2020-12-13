package ctl

import (
	"context"

	"github.com/eloylp/aton/internal/proto"
	"github.com/eloylp/aton/internal/video"
)

type Capturer interface {
	Start()
	Output() <-chan *video.Capture
	Close()
	Status() string
	UUID() string
}

type DetectorClient interface {
	Connect() error
	LoadCategories(ctx context.Context, request *proto.LoadCategoriesRequest) (*proto.LoadCategoriesResponse, error)
	Recognize(ctx context.Context) (chan<- *proto.RecognizeRequest, <-chan *proto.RecognizeResponse, error)
	Shutdown()
}
