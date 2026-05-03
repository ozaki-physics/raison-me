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

# 以下が スクリプト 本体

# 厳格モード
# -e はコマンドが失敗したらスクリプトを即終了
# -u は未設定(または空)の変数参照をエラー扱い
# -o pipefail は パイプライン内のどれかのコマンドが失敗した場合にパイプ全体を失敗扱い
set -euo pipefail

MIG_DIR="${MIGRATION_SQL_DIR:-migrations}"
OUT_DIR="${MIGRATION_OUTPUT_DIR:-release}"
# -p オプションで 親ディレクトリがなければ再帰的に作成
# ディレクトリ名が ハイフン で始まると オプションと誤認識される可能性があるため -- を付与
mkdir -p -- "$OUT_DIR"

# 直近までの適用を手元で把握するための「ローカルのschema_migrations」を信頼源にする想定。
# もしローカルDBがない運用なら、Gitのタグ/最新番号で代用して「全部」を束ねてもOK。
# ここでは「全ファイル」を束ねる簡易版を示す(本番はidempotent前提)。

ts=$(date +%F_%H%M%S)
out_file_path="$OUT_DIR/bundle_${ts}.sql"

# 現存のファイルから ハッシュ一覧を作る(一時ファイルを作成して 後で削除)
tmphash=$(mktemp)
for f in $(ls -1 "$MIG_DIR"/*.sql | sort); do
  base=$(basename "$f")
  # ファイルのハッシュを計算して tmphash に追記
  _sum=$(get_sha256 "$f")
  printf "$base $_sum\n" >> "$tmphash"
done

# bundle して 1個のファイルに書き出す
{
  printf '%s\n' "-- bundle generated at $ts"
  printf "BEGIN;\n"

  # schema_migrations(冪等)
  cat <<'SQL'
CREATE TABLE IF NOT EXISTS public.schema_migrations (
  id         bigserial PRIMARY KEY,
  filename   text NOT NULL UNIQUE,
  sha256     text NOT NULL,
  applied_at timestamptz NOT NULL DEFAULT now()
);
SQL

  # 各マイグレーションを順に（ファイル内の先頭/末尾の単独の BEGIN/COMMIT 行は取り除く）
  for f in $(ls -1 "$MIG_DIR"/*.sql | sort); do
    base=$(basename "$f")
    printf "\n"
    printf '%s\n' "-- >>> START $base >>>"
    # ファイルを行単位で読み込み、行全体が BEGIN; または COMMIT;（前後に空白許容）の場合はスキップ
    # それ以外の行はそのまま出力する
    awk 'BEGIN{IGNORECASE=1}
    /^[[:space:]]*BEGIN;[[:space:]]*$/ {next}
    /^[[:space:]]*COMMIT;[[:space:]]*$/ {next}
    {print}
    ' "$f"
    printf '%s\n' "-- <<< END $base <<<"
  done

  # ハッシュを一括登録(既にある filename はスキップ) + 改変検知(違うハッシュはエラー)
  printf "\n"
  # 匿名 PL/pgSQL ブロック で ブロック内部に シングルクォートがあっても エスケープしなくて済むように $$ で囲む

  # ダブルクオート だと $$ が 展開されてしまうため シングルクオートで囲む
  printf 'DO $$ BEGIN\n'
  # 現存ファイルから 計算したハッシュ一覧(tmphash)の各行を読み込み
  while read -r fname hash; do
    # 同じファイル名なのに ハッシュ値が違う場合は エラー
    printf "  IF EXISTS (SELECT 1 FROM public.schema_migrations WHERE filename = '$fname' AND sha256 <> '$hash') THEN\n"
    printf "    RAISE EXCEPTION 'Hash mismatch for $fname';\n"
    printf "  END IF;\n"
    # ハッシュを登録 (ON CONFLICT で filename が既にあれば何もしない)
    printf "  INSERT INTO public.schema_migrations(filename, sha256)\n"
    printf "  VALUES ('$fname', '$hash')\n"
    printf "  ON CONFLICT (filename) DO NOTHING;\n"
  done < "$tmphash"
  printf 'END; $$;\n'

  printf "COMMIT;\n"
} > "$out_file_path"

rm -f "$tmphash"

printf "✅ Created $out_file_path\n"
