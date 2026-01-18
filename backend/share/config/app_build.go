package config

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App interface {
	// グローバル設定
	GetConfig() Config
	// アプリとしての コンテキスト
	GetAppCtx() context.Context
	// DB の プール
	GetPool() *pgxpool.Pool
}

func NewAppBuild() App {
	log.Println("AppBuild: called")

	app := &app{}
	// 環境変数 など 読み込み
	app.config = NewConfig()

	// アプリとしての Context (アプリ全体で共有)
	appCtx := context.Background()
	app.appContext = appCtx

	// DB プール の 作成
	// appCtx を 引き継いで タイムアウト付き Context を 作成
	initCtx, cancel := context.WithTimeout(appCtx, 10*time.Second)
	defer cancel()
	pool, err := newPool(initCtx, app.config)
	if err != nil {
		log.Fatalf("Failed to create pool: %v", err)
	}
	app.pool = pool

	log.Println("AppBuild: completed")
	return app
}

// アプリケーション 全体 の 構造体
type app struct {
	config     Config
	appContext context.Context
	pool       *pgxpool.Pool
}

func (a *app) GetConfig() Config {
	return a.config
}

func (a *app) GetAppCtx() context.Context {
	return a.appContext
}

func (a *app) GetPool() *pgxpool.Pool {
	return a.pool
}
