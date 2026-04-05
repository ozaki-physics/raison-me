# DB 接続設計
PostgreSQL には `pgx` を直接使い、アプリ内の接続管理は `pgxpool` で行う。  
リクエストごとに `sql.Open()` はせず、接続プールをプロセス単位で初期化して再利用する。  
クエリ実行時は必ず上位から渡された `context.Context` を使い、必要に応じて timeout を付与する。

## 接続プール設計
DB pool はコネクションの再利用の仕組みであり、毎リクエストごとに張り直すものではない。

### `database/sql` を使う場合の考え方
- `sql.Open()` をリクエストごとに呼ばない
- `database/sql` が内部で DB pool を管理する

### 接続数
- DB に同時接続できる最大数は `db.SetMaxOpenConns(10)` のように設定する
- DB 側の `max_connections` と必ず整合させる
- Cloud Run では `インスタンス数 × MaxOpenConns = 実接続数` になる

### アイドル接続数
- 最大アイドル数は `db.SetMaxIdleConns(5)` のように設定する
- 少なすぎると毎回再接続になる
- 多すぎると DB 資源を無駄に占有する

### 接続寿命
- 接続の寿命は `db.SetConnMaxLifetime(30 * time.Minute)` のように設定する
- DB が切る前に Go 側から捨てる運用にする
- 切れたコネクションを再利用すると原因が分かりにくいエラーになりやすい

## クエリ実行時の扱い
`database/sql` を使う場合、`rows, err := db.QueryContext(ctx, query)` はプールから空いているコネクションを取得して使う。

- `rows.Close()` でコネクションをプールに返却する
- `rows.Close()` や `tx.Rollback()` を忘れるとリークする
- `tx, _ := db.BeginTx(ctx, nil)` は 1 本のコネクションを占有する

PgBouncer / Supavisor は Supabase 側に用意された pooler で、複数アプリ間の DB 接続を集約する。  
Go アプリ内で行う pool は、アプリ内 goroutine 間での接続利用を制御するものであり役割が異なる。

## 採用ライブラリ
PostgreSQL を使う前提なので `pgx` を直接使う。  
将来的に NewSQL を使う場合でも PostgreSQL 互換の選択肢が多いため、この方針を維持しやすい。  
PostgreSQL が 好きだし NewSQL にするとしても PostgreSQL 互換もあるから __pgx を直接使おう__  

### `pgx` の位置づけ
- `pgx` は PostgreSQL 専用に最適化された Go ネイティブドライバ
- `database/sql` インタフェース互換も含まれている
- `database/sql` は複数 DB を共通化するために あえて機能を削って意図的に抽象化されている

参考:
- https://github.com/jackc/pgx
- https://pkg.go.dev/github.com/jackc/pgx/v5#section-readme
- ライセンスは MIT

>pgx を `database/sql` 互換のドライバとして使用するには `github.com/jackc/pgx/v5/stdlib` をご利用ください。  
>詳細は当該パッケージのドキュメントを参照してください。  
>`*pgx.Conn` はデータベースへの単一接続を表し 並行処理安全ではありません  
>並行処理安全な接続プールには `github.com/jackc/pgx/v5/pgxpool` パッケージを使用してください。  

参考:
- https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool

>パッケージ pgxpool は、pgx 用の並行処理安全な接続プールです。  
>pgxpool は、pgx 接続とほぼ同一のインターフェースを実装しています。  

初期化例:

```go
pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
```

## `pgxpool` の使い方
SQL 実行時は返り値の形に応じて API を使い分ける。

- `pool.Exec()`: 返り値(行)が不要な更新系
- `pool.QueryRow()`: 1 行だけ欲しい
- `pool.Query()`: 複数行を読む

### `Query()` の注意点
- `Query()` は行のストリームを返す = サーバーから行を受け取り続ける状態 = ちゃんと close しないと コネクションがプールに戻らない
- `pool.Query()` が返す `rows` は単なるスライスではなく、DB からのストリーミング読み取り状態 -> `pgxpool` はそのコネクションを使用中扱いにする
- `rows.Close()` はストリームを破棄してコネクションをプールに戻す
- `QueryRow()` はストリーム管理が不要

### トランザクション
- pool を使う場合でも必ず Context を期限付きにする
- トランザクションは `pool.BeginTx()` を使う
- 実行は `tx.Exec` / `tx.QueryRow` / `tx.Query` を使う

## Context 設計
下層で勝手に `context.Background()` を作らない。  
なぜなら 上位がキャンセルしても止まらないクエリが生まれてしまう  
実装では Tx 開始 から 生まないため、上位から受け取った `ctx` をそのまま下へ流す。

参考:
- [ozaki-physics/go-training-composition: context パッケージ の勉強](https://github.com/ozaki-physics/go-training-composition/blob/develop/docs/package_context.md)

基本方針:
- Tx 開始から commit / rollback まで同じ `ctx` を使う
- ユースケース層で `context.WithTimeout(ctx, 2*time.Second)` を作ってリポジトリへ渡す
- リポジトリでは受け取った `ctx` をそのまま使う

例:

```go
row := r.pool.QueryRow(ctx, `select id, name from users where id=$1`, id)
```

## アプリ全体で使う Context
アプリ全体では主に 2 種類の Context を使い分ける。
1. App 自体の Context
- レスポンス作成に不要な処理は appCtx 側で行う
- 必ず timeout を付ける

2. HTTP ハンドラ配下の Context
- `r.Context()` を親にする
- ハンドラから下 (`UseCase -> Repo -> DB`) へは基本 `r.Context()` を流す
- リクエストの結果を作るための作業だから
- 必要なら短めの timeout を上書きする

プレゼン層で Context に対して行うこと:
- タイムアウトを付与する
- リクエストスコープ (`RequestID`, `userID`, `logger`, `tenantID`) を付与する

整理:
- `appCtx`: アプリ全体の親。停止シグナルで cancel される
- `initCtx`: `context.WithTimeout(appCtx, X秒)` (pool 初期化用: pool 初期化用の ctx は `appCtx` の子にするのがよい)
- `reqCtx`: `r.Context()` (クエリ実行用)

## Timeout API の使い分け
- `WithTimeout()`: 止まれば十分なときに使う, Go 1.7 (2016年8月)
- `WithTimeoutCause()`: 止まる理由を後で使いたいときに使う, Go 1.20 (2023年2月)

通常は `WithTimeout()` を使う  
理由が必要な箇所だけ `WithTimeoutCause()` を使う

## 関連メモ
### Go のメソッドドキュメント
参考:
- [ozaki-physics/go-training-composition: Go でのドキュメントの書き方](https://github.com/ozaki-physics/go-training-composition/blob/develop/docs/godoc_memo.md)

`// Deprecated: 非推奨を表す` のように書く。
すると以下のような表示になる
- [database/sql/driver ColumnConverter](https://pkg.go.dev/database/sql/driver@go1.25.6#ColumnConverter)

### Go のポインタ
参考:
- [ozaki-physics/go-training-composition: Go の勉強](https://github.com/ozaki-physics/go-training-composition/blob/main/docs/go_tour.md#ポインタ)
