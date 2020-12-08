package www_test

import (
	"net/http"
	"testing"

	"github.com/eloylp/kit/test/check"
	"github.com/eloylp/kit/test/handler"

	"github.com/eloylp/aton/internal/ctl/www"
)

func TestHandlers(t *testing.T) {
	cases := []handler.Case{
		{
			Case:     "Status is showing correctly",
			Path:     "/statuss",
			Method:   http.MethodGet,
			Checkers: []check.Function{check.HasStatus(http.StatusOK), check.ContainsJSON(`{"status":"ok"}`)},
		},
	}
	t.Run("Running handler tests ...", handler.Tester(cases, www.Router(), nil, nil))
}
