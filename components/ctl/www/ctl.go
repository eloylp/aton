package www

import (
	"context"
)

type Ctl interface {
	AddNode(addr string) (string, error)
	AddCapturer(ctx context.Context, uuid, url string) error
	LoadCategories(ctx context.Context, categories []string, image []byte) error
	Shutdown(ctx context.Context) error
}
