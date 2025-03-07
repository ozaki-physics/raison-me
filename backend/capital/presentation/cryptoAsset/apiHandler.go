package cryptoasset

import (
	"fmt"
	"net/http"
)

type Handler interface {
	Handler(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	apiController ApiController
}

func CreateApiHandler(api ApiController) Handler {
	return &handler{api}
}

func (h *handler) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		fmt.Fprintln(w, "\nGET だよ")
		err := h.apiController.Get(w, r)
		if err != nil {
			http.Error(w, "Internal Server Error\n"+err.Error(), 500)
		}
	case http.MethodPost:
		fmt.Fprintln(w, "\nPOST だよ")
		err := h.apiController.Post(w, r)
		if err != nil {
			http.Error(w, "Internal Server Error\n"+err.Error(), 500)
		}
	// case http.MethodPut:
	// 	fmt.Fprintln(w, "\nPUT だよ")
	// 	err := h.apiController.Put(w, r)
	// 	if err != nil {
	// 		http.Error(w, "Internal Server Error\n"+err.Error(), 500)
	// 	}
	// case http.MethodDelete:
	// 	fmt.Fprintln(w, "\nDELETE だよ")
	// 	err := h.apiController.Delete(w, r)
	// 	if err != nil {
	// 		http.Error(w, "Internal Server Error\n"+err.Error(), 500)
	// 	}
	default:
		fmt.Fprintln(w, "許可されたメソッドじゃないよ")
	}
}
