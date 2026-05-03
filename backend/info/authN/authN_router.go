// Context : 認証(Authentication: AuthN)
package authn

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ozaki-physics/raison-me/info/authN/infra"
	"github.com/ozaki-physics/raison-me/info/authN/presen"
	"github.com/ozaki-physics/raison-me/info/authN/usecase"
	"github.com/ozaki-physics/raison-me/share"
	"github.com/ozaki-physics/raison-me/share/config"
	sharemiddleware "github.com/ozaki-physics/raison-me/share/middleware"
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
	refreshTokenRepo, _ := infra.NewRefreshTokenRepoSQL(app.GetPool())
	passwordHasher, _ := infra.NewPasswordBcrypt(app.GetConfig().GetAuthNPepper())
	accessTokenService, _ := infra.NewAccessTokenHMAC(app.GetConfig().GetAuthNJWTSecret())
	refreshTokenGenerator := infra.NewRefreshTokenRandom()
	refreshTokenHasher, _ := infra.NewRefreshTokenSHA256(app.GetConfig().GetAuthNRefreshTokenPepper())

	authN, err := usecase.NewAuthN(
		userRepo,
		passRepo,
		refreshTokenRepo,
		passwordHasher,
		accessTokenService,
		refreshTokenGenerator,
		refreshTokenHasher,
	)
	if err != nil {
		log.Fatalf("authN の ユースケース 初期化に失敗しました: %v", err)
	}
	api := presen.NewAPICase(authN, app.GetConfig().IsLive())

	// 認証 が 不要な エンドポイント
	r.Group(func(r chi.Router) {
		r.Post("/signin", share.NewApiHandler(api.SignIn).Handler)
		r.Post("/signout", share.NewApiHandler(api.SignOut).Handler)
		r.Route("/token", func(r chi.Router) {
			r.Post("/refresh", share.NewApiHandler(api.RefreshAccessToken).Handler)
		})
		r.Post("/signup", share.NewApiHandler(api.CreateUser).Handler)
	})

	// 認証 が 必要な エンドポイント
	r.Group(func(r chi.Router) {
		r.Use(sharemiddleware.AuthN(app.GetConfig(), accessTokenService))
		r.Get("/me", share.NewApiHandler(api.GetMe).Handler)
		r.Route("/users", func(r chi.Router) {
			r.Get("/", share.NewApiHandler(api.GetUserList).Handler)
			r.Get("/{userID}", share.NewApiHandler(
				func(w http.ResponseWriter, req *http.Request) error {
					userID := chi.URLParam(req, "userID")
					return api.SearchUser(w, req, userID)
				},
			).Handler)
		})
	})

	// 機能ごとで バージョニング するため
	v1 := chi.NewRouter()
	v1.Mount("/v1", r)
	return v1
}
