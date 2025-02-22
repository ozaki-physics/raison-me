// Context : 認証(Authentication: AuthN)
package authn

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// AuthN コンテキスト を統括するルータ
func Router() chi.Router {
	// ルーティングの定義
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("これは info の authn だよ\n"))
	})

	return r
}
