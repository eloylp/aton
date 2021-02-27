package ctl

import (
	"context"
)

type NodeClient interface {
	Connect() error
	LoadCategories(ctx context.Context, r *LoadCategoriesRequest) error
	AddCapturer(ctx context.Context, r *AddCapturerRequest) error
	RemoveCapturer(ctx context.Context, r *RemoveCapturerRequest) error
	NextStatus() (*Status, error)
	NextResult() (*Result, error)
	Shutdown() error
}
