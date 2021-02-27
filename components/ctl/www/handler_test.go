package www_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/eloylp/kit/test/check"
	"github.com/eloylp/kit/test/handler"

	"github.com/eloylp/aton/components/ctl/metrics"
	"github.com/eloylp/aton/components/ctl/www"
)

var metricsService = metrics.NewService()

func TestHandlers(t *testing.T) {
	cases := []handler.Case{
		{
			Case:     "Status is showing correctly",
			Path:     "/status",
			Method:   http.MethodGet,
			Checkers: []check.Function{check.HasStatus(http.StatusOK), check.ContainsJSON(`{"status":"ok"}`)},
		},
		{
			Case:     "Metrics is showing correctly",
			Path:     "/metrics",
			Method:   http.MethodGet,
			Checkers: []check.Function{check.HasStatus(http.StatusOK), check.Contains(`aton_ctl`)},
		},
	}
	t.Run("Running handler tests ...", handler.Tester(cases, www.Router(NewFakeCtl(), metricsService.HTTPHandler()), func(t *testing.T) {
		metricsService.NodeUP("A1234")
	}, func(t *testing.T) {
		metricsService.NodeDown("A1234")
	}))
}

type FakeCtl struct {
}

func (f *FakeCtl) Shutdown(ctx context.Context) error {
	panic("implement me")
}

func NewFakeCtl() *FakeCtl {
	return &FakeCtl{}
}

func (f *FakeCtl) AddNode(addr string) (string, error) {
	panic("implement me")
}

func (f *FakeCtl) AddCapturer(ctx context.Context, uuid, url string) error {
	panic("implement me")
}
