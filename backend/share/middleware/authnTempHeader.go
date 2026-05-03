package middleware

import (
	"log"
	"net/http"

	"github.com/ozaki-physics/raison-me/share/config"
)

// 正規の認証 とは 別の ミドルウェア
func CustomMiddleware(config config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// パブリック な パス なら 認証なしで 通す
			if isPublicPath(r.URL.Path) {
				log.Printf("静的ファイル配信のため 認証スキップ: %s", r.URL.Path)
				next.ServeHTTP(w, r)
				return
			}

			// 独自ヘッダー のチェック
			if config.IsCloud() == true {
				customTempHeader := r.Header.Get("X-Custom-Temp-Header")
				isVerified := customTempHeader == config.GetSampleAPIToken()
				if !isVerified {
					log.Printf("カスタムヘッダーの一致: %v", isVerified)
					http.Error(w, "not access", http.StatusUnauthorized)
					return
				}
			} else {
				log.Printf("独自ヘッダー チェックをスキップ: ローカル環境のため")
			}

			next.ServeHTTP(w, r)
		})
	}
}
