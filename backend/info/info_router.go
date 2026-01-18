// Service : info のコンテキストたちの DI などを行う
package info

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	authn "github.com/ozaki-physics/raison-me/info/authN"
	"github.com/ozaki-physics/raison-me/share/config"
)

// info サービス を統括するルータ
// アプリ 全体共通の config.App を 受け取っているが サービス 固有の App を作ってもよい
func Router(app config.App) chi.Router {
	// ルーティングの定義
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("これは info だよ\n"))
	})

	r.Mount("/auth-n", authn.Router(app))
	return r
}
