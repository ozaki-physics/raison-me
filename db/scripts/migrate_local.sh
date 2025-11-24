#!/usr/bin/env bash

# ファイルの SHA-256 を計算して出力する メソッド
# sha256sum または openssl を試す
get_sha256(){
  local file="$1"
  local _sum
  # 環境に sha256sum コマンドが インストールされているか 確認
  if command -v sha256sum >/dev/null 2>&1; then
    # 出力: "<hash>  <filename>" なので 最初のフィールド を _sum に 格納
    _sum=$(sha256sum "$file" | awk '{print $1}')
  else
    # 出力: "SHA256(filename)= <hash>" か "<hash>" なので 最後のフィールド を _sum に 格納
    # $NF は awk の特殊変数で 最後のフィールド を意味する
    _sum=$(openssl dgst -sha256 "$file" | awk '{print $NF}')
  fi

  # _sum に値が入っていれば出力 なければ エラー
  if [ -n "${_sum:-}" ]; then
    printf '%s' "$_sum"
    return 0
  else
    printf 'no suitable hash command found\n' >&2
    return 1
  fi
}

# マイグレーションファイルの実行 と ハッシュの記録 を行う メソッド
apply(){
  local f="$1"
  printf '>> %s\n' "$f"
  # ファイルを実行
  psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -f "$f"
  # ハッシュを取得
  local sum
  sum=$(get_sha256 "$f")

  # マイグレーション履歴に 挿入
  # 接続文字列を -d で渡すのを 推奨
  # なぜなら -v を 接続文字列 の後ろに 書くと 変数展開の影響で 正しく解釈されない場合があるため
  schema_migrations_insert=$(
    psql -v ON_ERROR_STOP=1 -d "$DATABASE_URL" \
      -c "INSERT INTO public.schema_migrations(filename, sha256)
          VALUES ('$(basename "$f")', '$sum')
          ON CONFLICT (filename) DO NOTHING;"
  )
  printf '<< schema_migrations: %s\n' "$schema_migrations_insert"
}

# 以下が スクリプト 本体

# 厳格モード
# -e はコマンドが失敗したらスクリプトを即終了
# -u は未設定(または空)の変数参照をエラー扱い
# -o pipefail は パイプライン内のどれかのコマンドが失敗した場合にパイプ全体を失敗扱い
set -euo pipefail

# 環境変数 DATABASE_URL の存在 かつ 非空 をチェック
# DATABASE_URL が 未設定または空 だと スクリプトは エラーメッセージ(set DATABASE_URL for local) を出して終了
# DATABASE_URL の 想定は postgres://<ユーザー名>:<パスワード>@<ホスト>:<ポート>/<データベース名>
: "${DATABASE_URL:?set DATABASE_URL for local}"

# 空文字を許容したい場合は ${VAR:-default}
MIG_DIR="${MIGRATION_SQL_DIR:-migrations}"

# 初回用(冪等)
# マイグレーション履歴を保存するためのテーブル public.schema_migrations を作成
# 何度実行しても既存テーブルがあれば上書きせず安全にスキップ = 冪等性
# -v ON_ERROR_STOP=1 により psql は SQL 実行中にエラーが発生した時点で即座に終了コードを返す
psql "$DATABASE_URL" -v ON_ERROR_STOP=1 <<'SQL'
CREATE TABLE IF NOT EXISTS public.schema_migrations (
  id         bigserial PRIMARY KEY,
  filename   text NOT NULL UNIQUE,
  sha256     text NOT NULL,
  applied_at timestamptz NOT NULL DEFAULT now()
);
SQL

# マイグレーションディレクトリ内の SQL ファイルを 順に処理
for f in $(ls -1 "$MIG_DIR"/*.sql | sort); do
  # ファイル名だけ取得
  base=$(basename "$f")
  # -A はアンアラインド(ヘッダや余分な空白を省く)
  # -t は tuples-only(列名や行数の表示を抑制)
  # psql コマンドが失敗しても set -e の影響で スクリプト全体が終了しないように || true を付与
  # exists は 1行だけの出力 (存在すれば 1, 存在しなければ空文字)
  exists=$(psql -d "$DATABASE_URL" -At \
    -c "SELECT 1 FROM public.schema_migrations WHERE filename='$base' LIMIT 1" || true
  )

  # exists が空(未適用)なら apply 関数を呼んでファイルを実行
  if [ -z "$exists" ]; then
    apply "$f"
  else
    sum=$(get_sha256 "$f")

    dbsum=$(psql -d "$DATABASE_URL" -At \
      -c "SELECT sha256 FROM public.schema_migrations WHERE filename='$base'"
    )
    # 適用済み で 改変されているか検知
    if [ "$sum" != "$dbsum" ]; then
      printf '!! %s\n' "Hash mismatch for $base" >&2
      exit 1
    fi
  fi
done

printf "✅ Local migration done.\n"
