// Context : 認証(Authentication: AuthN)
package authn

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ozaki-physics/raison-me/info/authN/infra"
	"github.com/ozaki-physics/raison-me/info/authN/presen"
	"github.com/ozaki-physics/raison-me/info/authN/usecase"
	"github.com/ozaki-physics/raison-me/share"
)

// AuthN コンテキスト を統括するルータ
func Router() chi.Router {
	// ルーティングの定義
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("これは info の authn だよ\n"))
	})

	r.Get("/user", SearchUser)

	// DI のイメージ?
	// TODO: 用途不明
	userRepo, _ := infra.NewUserRepoJSON()
	credentialRepo, _ := infra.NewCredentialRepoJSON()
	cryptoRepo, _ := infra.NewCryptoRepoJSON()
	userUsecase, _ := usecase.NewUser(userRepo, cryptoRepo)
	credentialUsecase, _ := usecase.NewCredential(credentialRepo, cryptoRepo)
	api := presen.NewAPICase(userUsecase, credentialUsecase)

	// TODO: 用途不明
	r.Route("/entry", func(r chi.Router) {
		r.Get("/", share.NewApiHandler(api.UserCreate).Handler)
		r.Post("/", share.NewApiHandler(api.UserEntry).Handler)
	})
	r.Route("/token", func(r chi.Router) {
		r.Post("/", share.NewApiHandler(api.IDTokenGenerate).Handler)
	})

	return r
}

// TODO: マジで暫定的 な プレゼンテーション層
func SearchUser(w http.ResponseWriter, r *http.Request) {
	storagePath := infra.NewStoragePath()

	userCount := storagePath.GetUser()
	w.Write([]byte("User count: " + strconv.Itoa(userCount) + "\n"))
}
