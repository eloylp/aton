package www

import (
	"context"
)

type Ctl interface {
	AddDetector(addr string) (string, error)
	AddCapturer(ctx context.Context, uuid, url string) error
	Shutdown(ctx context.Context) error
}
