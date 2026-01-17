package middleware_temp

import (
	"log"
	"net/http"
	"strings"

	"github.com/ozaki-physics/raison-me/share/config"
)

// 認証 ミドルウェア の サンプル
func SampleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: 認証処理 を ここに 書く
		// リクエスト ヘッダー の Authorization から 値を取得して チェック
		requestAuthorization := r.Header.Get("Authorization")
		log.Printf("Authorization: %s", requestAuthorization)
		if requestAuthorization == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}

		token, ok := parseBearer(requestAuthorization)
		if !ok {
			http.Error(w, "invalid Authorization header format(expected 'Bearer <token>')", http.StatusUnauthorized)
			return
		}

		sampleAPIToken := config.NewConfig().GetSampleAPIToken()
		if token != sampleAPIToken {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
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
