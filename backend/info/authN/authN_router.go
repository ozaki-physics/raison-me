// Context : 認証(Authentication: AuthN)
package authn

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ozaki-physics/raison-me/info/authN/infra"
	"github.com/ozaki-physics/raison-me/info/authN/presen"
	"github.com/ozaki-physics/raison-me/info/authN/usecase"
	"github.com/ozaki-physics/raison-me/share"
	"github.com/ozaki-physics/raison-me/share/config"
)

// AuthN コンテキスト を統括するルータ
// アプリ 全体共通の config.App を 受け取っているが サービス 固有の App を作ってもよい
func Router(app config.App) chi.Router {
	// ルーティングの定義
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("これは info の authn だよ\n"))
	})

	// DI
	userRepo, _ := infra.NewUserRepoSQL(app.GetPool())
	passRepo, _ := infra.NewPassRepoSQL(app.GetPool())
	passwordHasher, _ := infra.NewPasswordBcrypt(app.GetConfig().GetAuthNPepper())
	authN, _ := usecase.NewAuthN(userRepo, passRepo, passwordHasher)
	api := presen.NewAPICase(authN)

	r.Post("/signin", share.NewApiHandler(api.SignIn).Handler)
	r.Get("/signout", share.NewApiHandler(api.SignOut).Handler)
	r.Get("/me", share.NewApiHandler(api.GetMe).Handler)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", share.NewApiHandler(api.CreateUser).Handler)
		r.Get("/{userID}", share.NewApiHandler(
			func(w http.ResponseWriter, req *http.Request) error {
				userID := chi.URLParam(req, "userID")
				return api.SearchUser(w, req, userID)
			},
		).Handler)
	})

	r.Route("/token", func(r chi.Router) {
		r.Post("/", share.NewApiHandler(api.GenerateToken).Handler)
		r.Get("/verify", share.NewApiHandler(api.VerifyToken).Handler)
		r.Post("/refresh", share.NewApiHandler(api.RefreshToken).Handler)
	})

	// TODO: 暫定
	r.Get("/users-all", share.NewApiHandler(api.GetUserList).Handler)

	// 機能ごとで バージョニング するため
	v1 := chi.NewRouter()
	v1.Mount("/v1", r)
	return v1
}
