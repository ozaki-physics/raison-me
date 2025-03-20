package share

import (
	"net/http"
)

type Handler interface {
	// ServeHTTP(w http.ResponseWriter, r *http.Request)
	Handler(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	ac apiCase
}

type apiCase func(w http.ResponseWriter, r *http.Request) error

func NewApiHandler(ac apiCase) Handler {
	return &handler{ac}
}

func (h *handler) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := h.ac(w, r)
	if err != nil {
		http.Error(w, "Internal Server Error\n"+err.Error(), 500)
	}
}
