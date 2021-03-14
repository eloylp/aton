package www

import (
	"encoding/json"
	"fmt"
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

func AddNodeHandler(ctl Ctl) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &AddNodeRequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			reply(w, err.Error(), http.StatusBadRequest)
			return
		}
		uuid, err := ctl.AddNode(req.Addr)
		if err != nil {
			reply(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp := &AddNodeResponse{UUID: uuid}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			reply(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func AddTargetHandler(ctl Ctl) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &AddTargetRequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			reply(w, err.Error(), http.StatusBadRequest)
			return
		}
		err := ctl.AddCapturer(r.Context(), req.UUID, req.TargetAddr)
		if err != nil {
			reply(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reply(w, fmt.Sprintf("added node %s", req.UUID), http.StatusOK)
	}
}

func LoadCategoriesHandler(ctl Ctl) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &LoadCategoriesRequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			reply(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := ctl.LoadCategories(r.Context(), req.Categories, req.Image); err != nil {
			reply(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reply(w, "categories loaded in all nodes", http.StatusOK)
	}
}

func reply(w http.ResponseWriter, msg string, status int) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(&struct {
		Message string `json:"message"`
	}{
		Message: msg,
	}); err != nil {
		panic(err)
	}
}
