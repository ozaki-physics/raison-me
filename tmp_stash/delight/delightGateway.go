// Service : delight のコンテキストたちの DI などを行う
package delight

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("これは delight だよ\n"))
	})

	return r
}
