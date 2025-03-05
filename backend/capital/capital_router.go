// Service : capital のコンテキストたちの DI などを行う
package capital

import (
	"net/http"

	"github.com/go-chi/chi"
)

// capital サービス を統括するルータ
func Router() chi.Router {
	// ルーティングの定義
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("これは capital だよ\n"))
	})

	return r
}
