package handler

import (
	"encoding/json"
	"net/http"
)

func StatusHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(struct {
			Status string `json:"status"`
		}{
			Status: "ok",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
