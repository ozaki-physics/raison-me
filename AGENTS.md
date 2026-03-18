# Repository Guidelines

## プロジェクト構成 と モジュール配置
このリポジトリは Go バックエンドと SvelteKit フロントエンドで構成されています。
`backend/` は `backend/main.go` から起動する API で、機能別ルータは `capital/`、`info/`、`growth/`、`seed/` などに配置されています。
共通設定や DB 接続は `backend/share/config/` にあります。
`frontend/raison-me-fe/` はフロントエンド本体で、画面ルートは `src/routes/`、共通コードは `src/lib/`、静的ファイルは `static/` に置きます。
設計メモや仕様は `docs/`、CI/CD は `.github/workflows/` を参照してください。

ADR を書くときは `docs/adr` の index を参照する
固まった設計 を書くときは `docs/design-docs` の index を参照する

## ビルド, テスト, 開発 コマンド
コマンドはリポジトリ直下ではなく、各アプリのディレクトリで実行してください。

- `cd backend && go run .` : API をローカル起動します。`PORT` 未指定時は `8081` です。
- `cd backend && go test ./...` : Go の全テストを実行します。
- `cd backend && go build -o main main.go` : バックエンドの実行バイナリを作成します。
- `cd frontend/raison-me-fe && npm install` : フロントエンド依存をインストールします。
- `cd frontend/raison-me-fe && npm run dev` : Vite 開発サーバを `8082` で起動します。
- `cd frontend/raison-me-fe && npm run build` : 本番ビルドを作成します。
- `cd frontend/raison-me-fe && npm run check && npm run lint` : 型検査、整形確認、Lint をまとめて実行します。

## コーディング規約と命名
Go は `gofmt` 前提で、公開識別子は PascalCase、非公開識別子は camelCase を使います。
フロントエンドは Prettier と ESLint を基準にし、`.prettierrc` ではタブ、シングルクォート、末尾カンマなし、`printWidth: 100` が設定されています。
Svelte のルートファイルは `+page.svelte` や `+layout.svelte` の命名を守ってください。
Go のファイル名は `authN_router.go`、`json_detail_test.go` のように、内容が分かる lower_snake_case を使います。

## コーディング思想
- ドメイン駆動開発(DDD)を行う
- DDD の アプリケーション層 に該当する概念は ユースケース層 と呼称する
- プレゼンテーション層は薄く保ち 業務ロジック は ユースケース層 と ドメイン層 に配置する
- ユースケース層 から プレゼン層 は Dto オブジェクト で連携する
- プレゼン層 から ユースケース層 は param オブジェクト で連携する
- 外部 から プレゼン層 は request オブジェクト で連携する
- プレゼン層 から 外部 は response オブジェクト で連携する
- 凝集度を高くする
- ユースケース層 で DTO に詰め替える意味は プレゼン層 で 業務ロジック を書かないようにするため(ドメインオブジェクト を プレゼン層 に露出すると プレゼン層 で 業務ロジック が書ける)
- DTO の 中では ドメインオブジェクト を含めない
- param は プレゼン層 から ユースケース層 に渡す情報が多くなったときに 使う(渡す情報が1,2個なら シンプルを優先し用意しない)

## テスト方針
バックエンドのテストは Go 標準の `testing` パッケージを使い、対象コードの近くに `*_test.go` として配置します。
バリデーション、JSON ヘルパー、ルーティングまわりを変更した場合は、対応するユニットテストを追加してください。
フロントエンドには専用のテストランナーがまだ入っていないため、少なくとも PR 前に `npm run check` と `npm run lint` を通してください。

## 設定ファイル

## コミットとプルリクエスト
ローカルフックは `git config core.hooksPath .githooks` で有効化できます。
`commit-msg` フックでは `feat:`、`fix:`、`docs:`、`refactor:`、`test:`、`chore:` の接頭辞が必須です。
履歴も同じ形式で揃っています。`main` と `master` への直接コミットは禁止です。
PR では変更概要を簡潔に書き、`Fixes #...` 形式で Issue を関連付けてください。
UI や API の挙動が変わる場合は、スクリーンショットやリクエスト例も添えてください。

## セキュリティと設定
実運用の秘密情報はコミットしないでください。
設定値は `backend/share/secrets/` と `frontend/raison-me-fe/.env.example` のサンプルを基準に管理し、新しい環境変数を追加した場合は PR に必要事項を明記してください。

## 禁止事項
`Set-Content` コマンドを __使わない__
なぜなら ファイルに差分が 無くても Git では 更新した扱いになってしまうから
