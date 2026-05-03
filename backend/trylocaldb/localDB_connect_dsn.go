package trylocaldb

import (
	"database/sql"
	"fmt"
	"log"

	globalConfig "github.com/ozaki-physics/raison-me/share/config"

	// 無いと言われたら go get で インストール すること
	// pwd: /app/backend
	// go get github.com/jackc/pgx/v5/stdlib
	// go mod tidy
	_ "github.com/jackc/pgx/v5/stdlib"
)

func getRecord() {
	// ホスト は Docker Compose のサービス名
	// postgres://<ユーザー名>:<パスワード>@<ホスト>:<ポート>/<データベース名>
	dsn := globalConfig.NewConfig().GetDSN()

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("sql.Open: %v\n", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("ping: %v\n", err)
	}

	var now string
	if err := db.QueryRow("SELECT now()::text").Scan(&now); err != nil {
		log.Fatalf("query: %v\n", err)
	}

	fmt.Printf("Connected to Postgres (pgx), time: %s\n", now)

	type Coin struct {
		ID              int
		Symbol          string
		CoinMarketCapID int
		CreatedAt       string
		UpdatedAt       string
		DeletedAt       sql.NullString
	}

	coins := []Coin{}
	rows, err := db.Query("SELECT * FROM app.capital_coin")
	if err != nil {
		log.Fatalf("query: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c Coin
		if err := rows.Scan(&c.ID, &c.Symbol, &c.CoinMarketCapID, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt); err != nil {
			log.Fatalf("scan: %v\n", err)
		}
		coins = append(coins, c)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("rows: %v\n", err)
	}
	defer rows.Close()
	fmt.Printf("Coins: %+v\n", coins)
}
