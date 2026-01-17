package trysupabase

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	globalConfig "github.com/ozaki-physics/raison-me/share/config"
)

// トランザクション プール での 接続確認
func transactionPooler_Connect() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	escapedPass := url.QueryEscape(config.Password)
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.User, escapedPass, config.Host, config.TransactionPort, config.DBName)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	ctx := context.Background()
	row := conn.QueryRow(ctx, "SELECT version()")

	var version string
	if err := row.Scan(&version); err != nil {
		if err == pgx.ErrNoRows {
			log.Println("No version row returned")
		} else {
			log.Fatalf("Failed to scan version: %v", err)
		}
		return
	}

	log.Println("PostgreSQL version:", version)
}

// セッション プール での 接続確認
func sessionPooler_Connect() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	escapedPass := url.QueryEscape(config.Password)
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.User, escapedPass, config.Host, config.SessionPort, config.DBName)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	sql := `SELECT version();`
	// Example query to test connection
	var version string
	err02 := conn.QueryRow(context.Background(), sql).Scan(&version)
	if err02 != nil {
		log.Fatalf("Query failed: %v", err02)
	}

	log.Println("Connected to:", version)
}

// トランザクション プール での ユーザー 一覧を取得
func transactionPooler_ReadUsers() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	escapedPass := url.QueryEscape(config.Password)
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.User, escapedPass, config.Host, config.TransactionPort, config.DBName)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	sql := `SELECT * from app.users;`
	type User struct {
		ID   int
		Name string
	}

	rows, err := conn.Query(context.Background(), sql)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	var results []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			log.Fatalf("Row scan failed: %v", err)
		}
		results = append(results, user)
	}

	if rows.Err() != nil {
		log.Fatalf("Rows iteration error: %v", rows.Err())
	}

	log.Println("Fetched users: ", results)
}

// セッション プール での ユーザー 一覧を取得
func sessionPooler_ReadUsers() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	escapedPass := url.QueryEscape(config.Password)
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.User, escapedPass, config.Host, config.SessionPort, config.DBName)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	sql := `SELECT * from app.users;`
	type User struct {
		ID   int
		Name string
	}

	rows, err := conn.Query(context.Background(), sql)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	var results []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			log.Fatalf("Row scan failed: %v", err)
		}
		results = append(results, user)
	}

	if rows.Err() != nil {
		log.Fatalf("Rows iteration error: %v", rows.Err())
	}

	log.Println("Connected to:", results)
}

// Pool を 作成
// チャピの通りに書いてみる
// key.json から 情報を 取得 して DSN を 作成
func newPool_keyJSON(ctx context.Context) (*pgxpool.Pool, error) {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	escapedPass := url.QueryEscape(config.Password)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=require", config.User, escapedPass, config.Host, config.TransactionPort, config.DBName)

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// 重要: Transaction pooler 対策（prepared statements を実質オフにする）
	// pgx はデフォルトで extended protocol + prepare/cache を使うため、PGBouncer系と相性が悪いことがある
	cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	// 接続数は Cloud Run の concurrency / max instances と合わせて設計（例は控えめ）
	cfg.MaxConns = 10
	cfg.MinConns = 0
	cfg.MaxConnIdleTime = 30 * time.Second
	cfg.MaxConnLifetime = 30 * time.Minute

	// タイムアウトなども好みで
	cfg.HealthCheckPeriod = 30 * time.Second

	return pgxpool.NewWithConfig(ctx, cfg)
}

// key.json から 作った Pool を 使う
// チャピの通りに書いてみる
func transactionPooler_ReadUsers_keyJSON() {
	ctx := context.Background()
	pool, err := newPool_keyJSON(ctx)
	if err != nil {
		log.Fatalf("Failed to create pool: %v", err)
	}
	defer pool.Close()

	sql := `SELECT * from app.users;`
	type User struct {
		ID   int
		Name string
	}

	rows, err := pool.Query(ctx, sql)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	var results []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			log.Fatalf("Row scan failed: %v", err)
		}
		results = append(results, user)
	}

	if rows.Err() != nil {
		log.Fatalf("Rows iteration error: %v", rows.Err())
	}

	log.Println("Fetched users: ", results)
}

// Pool を 作成
// チャピの通りに書いてみる
// globalConfig から 情報を 取得 して DSN を 作成
func newPool_globalConfig(ctx context.Context) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(globalConfig.NewConfig().GetSupabaseDSN())
	if err != nil {
		return nil, err
	}

	// 重要: Transaction pooler 対策（prepared statements を実質オフにする）
	// pgx はデフォルトで extended protocol + prepare/cache を使うため、PGBouncer系と相性が悪いことがある
	cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	// 接続数は Cloud Run の concurrency / max instances と合わせて設計（例は控えめ）
	cfg.MaxConns = 10
	cfg.MinConns = 0
	cfg.MaxConnIdleTime = 30 * time.Second
	cfg.MaxConnLifetime = 30 * time.Minute

	// タイムアウトなども好みで
	cfg.HealthCheckPeriod = 30 * time.Second

	return pgxpool.NewWithConfig(ctx, cfg)
}

// globalConfig から 作った Pool を 使う
// チャピの通りに書いてみる
func transactionPooler_ReadUsers_globalConfig() string {
	ctx := context.Background()
	pool, err := newPool_globalConfig(ctx)
	if err != nil {
		log.Fatalf("Failed to create pool: %v", err)
	}
	defer pool.Close()

	sql := `SELECT * from app.users;`
	type User struct {
		ID   int
		Name string
	}

	rows, err := pool.Query(ctx, sql)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	var results []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			log.Fatalf("Row scan failed: %v", err)
		}
		results = append(results, user)
	}

	if rows.Err() != nil {
		log.Fatalf("Rows iteration error: %v", rows.Err())
	}

	log.Println("Fetched users: ", results)
	return fmt.Sprintf("%v", results)
}
