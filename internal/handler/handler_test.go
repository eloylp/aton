package handler_test

import (
	"net/http"
	"testing"

	"github.com/eloylp/aton/internal/handler"
	"github.com/eloylp/kit/test/check"
	handlertest "github.com/eloylp/kit/test/handler"
)

func TestHandlers(t *testing.T) {
	cases := []handlertest.Case{
		{
			Case:     "Status is showing correctly",
			Path:     "/status",
			Method:   http.MethodGet,
			Checkers: []check.Function{check.HasStatus(http.StatusOK), check.ContainsJSON(`{"status":"ok"}`)},
		},
	}
	t.Run("Running handler tests ...", handlertest.Tester(cases, handler.Router(), nil, nil))
}
