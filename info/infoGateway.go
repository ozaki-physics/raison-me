// Service : info のコンテキストたちの DI などを行う
package info

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// info サービス を統括するルータ
func Router() chi.Router {
	// ルーティングの定義
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("これは info だよ\n"))
	})

	r.Mount("/authn", routerAuthN())
	return r
}

// AuthN コンテキスト を統括するルータ
func routerAuthN() chi.Router {
	// ルーティングの定義
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("これは info の authn だよ\n"))
	})

	return r
}
