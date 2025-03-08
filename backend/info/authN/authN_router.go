// Context : 認証(Authentication: AuthN)
package authn

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ozaki-physics/raison-me/info/authN/infra"
)

// AuthN コンテキスト を統括するルータ
func Router() chi.Router {
	// ルーティングの定義
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("これは info の authn だよ\n"))
	})

	r.Get("/user", SearchUser)
	return r
}

// TODO: マジで暫定的 な プレゼンテーション層
func SearchUser(w http.ResponseWriter, r *http.Request) {
	storagePath := infra.NewStoragePath()

	userCount := storagePath.GetUser()
	w.Write([]byte("User count: " + strconv.Itoa(userCount) + "\n"))
}
