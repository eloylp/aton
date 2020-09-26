package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Router() http.Handler {
	r := mux.NewRouter()
	r.Path("/status").Methods(http.MethodGet).HandlerFunc(StatusHandler())
	return r
}
