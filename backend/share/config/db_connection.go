package config

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool を 作成
// https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool
// TODO: もっと 適切な場所に 移動 する
func newPool(ctx context.Context, gc Config) (*pgxpool.Pool, error) {
	log.Println("Creating new DB pool")
	cfg, err := pgxpool.ParseConfig(gc.GetDSN())
	if err != nil {
		return nil, err
	}

	// Superbase の Transaction pooler や Cloud Run の 対応
	// Transaction pooler を主体にするため prepared statements を実質オフにする
	// pgx はデフォルトで extended protocol + prepare/cache を使うため PGBouncer 系と相性が悪いことがある
	// PGBouncer の場合は session pooling で prepared statements を使って パースを有効活用する方法もある
	// prepared statements とは SQL 文を事前にパースしておき 実行時にはパース済みのものを使う仕組み
	// SQL インジェクション対策 や パースのコスト削減 が 期待できる
	cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	// 接続数は Cloud Run の concurrency / max instances と合わせて設計
	cfg.MaxConns = 3
	cfg.MinConns = 0
	cfg.MaxConnIdleTime = 30 * time.Second
	cfg.MaxConnLifetime = 30 * time.Minute

	// タイムアウトなども好みで
	cfg.HealthCheckPeriod = 30 * time.Second

	// プール の 作成
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// Pool が 使えるか 一応 確認
	if err := pingPool(ctx, pool); err != nil {
		return nil, err
	}

	return pool, nil
}

// Pool が 使えるか確認
func pingPool(ctx context.Context, pool *pgxpool.Pool) error {
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return err
	}
	return nil
}
