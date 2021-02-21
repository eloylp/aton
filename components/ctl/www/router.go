package www

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Router(ctl Ctl, metrics http.Handler) http.Handler {
	r := mux.NewRouter()
	r.Path("/status").Methods(http.MethodGet).HandlerFunc(StatusHandler())
	r.Path("/metrics").Methods(http.MethodGet).Handler(metrics)
	return r
}
