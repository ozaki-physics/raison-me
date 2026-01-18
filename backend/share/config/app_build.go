package config

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App interface {
	// グローバル設定
	GetConfig() Config
	// DB の プール
	GetPool() *pgxpool.Pool
}

func NewAppBuild() App {
	log.Println("AppBuild: called")

	app := &app{}
	// 環境変数 など 読み込み
	app.Config = NewConfig()

	// DB プール の 作成
	// TODO: ここで Context 作っていいの?
	ctx := context.Background()
	pool, err := newPool(ctx, app.Config)
	if err != nil {
		log.Fatalf("Failed to create pool: %v", err)
	}
	app.Pool = pool

	log.Println("AppBuild: completed")
	return app
}

// アプリケーション 全体 の 構造体
type app struct {
	// グローバル設定
	Config Config
	// DB の プール
	Pool *pgxpool.Pool
}

func (a *app) GetConfig() Config {
	return a.Config
}

func (a *app) GetPool() *pgxpool.Pool {
	return a.Pool
}
