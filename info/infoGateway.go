// Service : info のコンテキストたちの DI などを行う
package info

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	authn "github.com/ozaki-physics/raison-me/info/authN"
)

// info サービス を統括するルータ
func Router() chi.Router {
	// ルーティングの定義
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("これは info だよ\n"))
	})

	r.Mount("/auth-n", authn.Router())
	return r
}
