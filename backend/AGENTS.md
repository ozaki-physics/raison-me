# Repository Guidelines

## プロジェクト構成とモジュール配置
`main.go` が HTTP サーバを起動し、`chi` で各機能ルータをマウントします。
`capital/`、`info/`、`growth/`、`seed/`、`regung/`、`zeit/` は機能単位のディレクトリです。
共通処理は `share/` にあり、特に `share/config/` で環境変数や DB 接続を扱います。
静的ファイルは `web/`、サンプル JSON や secret の雛形は `share/secrets/` と各機能配下の `json/` にあります。

ADR を書くときは 1個上の階層の `docs/adr` の index を参照する
固まった設計 を書くときは 1個上の階層の `docs/design-docs` の index を参照する

## ビルド・テスト・開発コマンド
作業ディレクトリは `backend/` を前提にします。

- `go run .` : ローカルサーバを起動します。`PORT` 未指定時は `8081` です。
- `go test ./...` : 全パッケージのユニットテストを実行します。
- `go test ./info/authN/...` : 認証まわりだけを重点的に確認します。
- `go build -o main main.go` : 本番用バイナリをビルドします。
- `docker build -t raison-me-backend .` : Cloud Run 向けイメージをローカルで作成します。

ローカル実行前に 必要な secret ファイルを用意してください。
参考: [`share/secrets/.env.example`](/mnt/c/Users/Owner/Documents/github_repository/raison-me/backend/share/secrets/.env.example)

## コーディング規約と命名
Go の標準に従い、整形は `gofmt` を前提にします。
インデントはタブ、公開識別子は PascalCase、非公開は camelCase を使います。
ファイル名は既存コードに合わせて lower_snake_case とし、例は `user_id_value.go`、`json_detail_test.go` です。
ルータ定義は `*_router.go`、層構造は `domain/`、`usecase/`、`infra` または `infrastructure/`、`presentation/` を踏襲してください。
`dto` は usecase 層 から presen 層 に渡すオブジェクト
ドメイン知識を presen 層で操作させないために dto の型は プリミティブ寄り にする

## テスト方針
テストは標準の `testing` パッケージを使い、対象ファイルの近くに `*_test.go` として配置します。
`info/authN/domain/` にあるようなテーブル駆動テストを基本にし、正常系と異常系の両方を確認してください。
値オブジェクトのバリデーションや永続化パスの分岐は優先してカバーします。

## コミットとプルリクエスト
Git 履歴では `feat:`、`fix:`、`docs:`、`chore:` の形式が使われています。
件名は短く、対象が分かる粒度で書いてください。例: `feat: DB pool を受け取るように修正`。
PR には変更した機能、必要な設定変更、関連 Issue や設計資料へのリンクを含め、API や HTML 出力を変える場合はリクエスト例や画面差分も添えてください。

## セキュリティと設定
実際の secret 値はコミットしないでください。
ローカルでは `./share/secrets/`、Cloud Run では `/app/share/secrets/` から設定を読み込みます。
設定項目を追加した場合はサンプルファイルも更新し、PR に必要な環境変数を明記してください。
