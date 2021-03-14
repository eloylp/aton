package www

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Router(ctl Ctl, metrics http.Handler) http.Handler {
	r := mux.NewRouter()

	// Application routes
	r.Path("/nodes").Methods(http.MethodPost).HandlerFunc(AddNodeHandler(ctl))
	r.Path("/targets").Methods(http.MethodPost).HandlerFunc(AddTargetHandler(ctl))
	r.Path("/categories").Methods(http.MethodPost).HandlerFunc(LoadCategoriesHandler(ctl))

	// Telemetry routes
	r.Path("/status").Methods(http.MethodGet).HandlerFunc(StatusHandler())
	r.Path("/metrics").Methods(http.MethodGet).Handler(metrics)

	return r
}
