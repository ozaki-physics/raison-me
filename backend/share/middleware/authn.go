package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	// TODO: share なのに サービス固有のドメインを import しているのは おかしいので どうにかする
	// share は interface だけにして 実装を サービス側に置くとか?
	// とりあえず middleware は presen 層 と考えられるため 問題ないとする
	"github.com/ozaki-physics/raison-me/info/authN/domain"
	"github.com/ozaki-physics/raison-me/info/authN/usecase"
	"github.com/ozaki-physics/raison-me/share/config"
)

type authContextKey string

const accessTokenClaimsContextKey authContextKey = "authN.accessTokenClaims"

func AuthN(config config.Config, accessTokenVerifier usecase.AccessTokenVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// パブリック な パス なら 認証なしで 通す
			if isPublicPath(r.URL.Path) {
				log.Printf("静的ファイル配信のため 認証スキップ: %s", r.URL.Path)
				next.ServeHTTP(w, r)
				return
			}

			// リクエスト ヘッダー の Authorization から 値を取得して チェック
			requestAuthorization := r.Header.Get("Authorization")
			if requestAuthorization == "" {
				http.Error(w, "missing Authorization header", http.StatusUnauthorized)
				return
			}

			// Authorization ヘッダー の 値を "Bearer <token>" 形式で パースする
			token, ok := parseBearer(requestAuthorization)
			if !ok {
				http.Error(w, "invalid Authorization header format(expected 'Bearer <token>')", http.StatusUnauthorized)
				return
			}

			// トークンを検証して クレームを取得する
			claims, err := accessTokenVerifier.Verify(token)
			if err != nil {
				log.Printf("アクセストークンの検証に失敗: %v", err)
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			// クレーム を コンテキスト に 保存して 次のハンドラーへ渡す
			ctx := WithAccessTokenClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// アクセストークンのクレーム を コンテキスト に 保存するための ヘルパー関数
func WithAccessTokenClaims(ctx context.Context, claims *domain.AccessTokenClaims) context.Context {
	return context.WithValue(ctx, accessTokenClaimsContextKey, claims)
}

// コンテキスト から アクセストークンのクレーム を 取得するための ヘルパー関数
func AccessTokenClaimsFromContext(ctx context.Context) (*domain.AccessTokenClaims, bool) {
	// TODO: 型アサーション で 変換してもいいのか?
	// ちゃんと New メソッド を 使ったほうがよい?
	claims, ok := ctx.Value(accessTokenClaimsContextKey).(*domain.AccessTokenClaims)
	return claims, ok
}

func parseBearer(auth string) (string, bool) {
	// "Bearer xxx" 以外は拒否
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 {
		return "", false
	}
	if !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", false
	}

	return token, true
}
