// Service : seed のコンテキストたちの DI などを行う
package seed

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("これは seed だよ\n"))
	})

	return r
}
