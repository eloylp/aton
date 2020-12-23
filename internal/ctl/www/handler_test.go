package www_test

import (
	"net/http"
	"testing"

	"github.com/eloylp/kit/test/check"
	"github.com/eloylp/kit/test/handler"

	"github.com/eloylp/aton/internal/ctl/metrics"
	"github.com/eloylp/aton/internal/ctl/www"
)

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
			Checkers: []check.Function{check.HasStatus(http.StatusOK), check.Contains(`aton_detector_up{uuid="A1234"} 1`)},
		},
	}
	t.Run("Running handler tests ...", handler.Tester(cases, www.Router(metrics.NewHTTPHandler()), func(t *testing.T) {
		metrics.NewRegister().DetectorUP("A1234")
	}, func(t *testing.T) {
		metrics.NewRegister().DetectorDown("A1234")
	}))
}
