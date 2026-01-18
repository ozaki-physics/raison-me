package middleware_temp

import (
	"log"
	"net/http"
	"strings"

	"github.com/ozaki-physics/raison-me/share/config"
)

// 認証 ミドルウェア の サンプル
func SampleMiddleware(config config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: 認証処理 を ここに 書く

			// 静的ファイル配信 のため ルート直下の以下のファイルへの アクセス は 認証不要 とする
			// TODO: 本来は 静的ファイル配信 自体を ミドルウェアの外に出すべき?
			// TODO: 静的ファイル も 認証が必要な場合は どうしよう
			staticFiles := []string{
				"/favicon.ico",
				"/robots.txt",
				"/sitemap.xml",
				"/humans.txt",
			}
			requestPath := r.URL.Path
			for _, file := range staticFiles {
				if requestPath == file {
					log.Printf("静的ファイル配信のため 認証スキップ: %s", requestPath)
					next.ServeHTTP(w, r)
					return
				}
			}

			// 開発環境 なら トークン認証なしで 通す
			if config.IsCloud() == false {
				log.Printf("開発環境と判定 認証をスキップ")
				next.ServeHTTP(w, r)
				return
			}

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

			sampleAPIToken := config.GetSampleAPIToken()
			if token != sampleAPIToken {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
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
