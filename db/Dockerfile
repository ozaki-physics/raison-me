# このコンテナの目的は ローカルで DB を 再現するため
# 本番環境は BigQuery や Supabase を検討する
# 開発に BigQuery や Supabase に直接接続していたら 費用がかかる
# かつ 開発, テスト, 本番 で データを分離するため
# Supabase と互換のあるバージョンとして 17 系を使う
# BigQuery と互換のあるバージョンは 不明
FROM postgres:17

# ビルド時に 引数 を 渡さなかった場合の デフォルト値
# ビルド時に 上書き可能
ARG POSTGRES_DB=dev_db
ARG POSTGRES_USER=dev_user

# 環境変数の設定
ENV POSTGRES_DB=${POSTGRES_DB}
ENV POSTGRES_USER=${POSTGRES_USER}

# 初期データの投入
# 公式イメージの 初回起動時にデータベースを初期化する機構を利用
# docker-entrypoint-initdb.d/ 配下に SQL ファイルを コピーすると コンテナ起動時に 自動で 実行される
COPY ./db/migrations/0000-init.sql /docker-entrypoint-initdb.d/

# PostgreSQL の デフォルト ポート を 開放
EXPOSE 5432

# ヘルスチェック の 設定
HEALTHCHECK \
  # 10秒ごとに 接続確認 を行う
  --interval=10s \
  # 応答が 5秒以内に返ってこなかったら タイムアウト とみなす
  --timeout=5s \
  # 開始から 10秒間は ヘルスチェックを開始しない
  --start-period=10s \
  # 失敗した場合は 最大 3回まで 再試行 する
  --retries=3 \
  # 接続可能なら 0 を返す, 不可能なら 非ゼロ値 を返す
  # pg_isready コマンド は 公式イメージに含まれている
  CMD pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB"
