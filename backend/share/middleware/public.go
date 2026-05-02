package middleware

// 静的ファイル配信 のため ルート直下の以下のファイルへの アクセス は 認証不要 とする
// TODO: 本来は 静的ファイル配信 自体を ミドルウェアの外に出すべき?
// TODO: 静的ファイル も 認証が必要な場合は どうしよう
func isPublicPath(requestPath string) bool {
	staticFiles := []string{
		"/",
		"/healthz",
		"/favicon.ico",
		"/robots.txt",
		"/sitemap.xml",
		"/humans.txt",
	}
	for _, file := range staticFiles {
		if requestPath == file {
			return true
		}
	}

	return false
}
