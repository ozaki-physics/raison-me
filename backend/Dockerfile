# ビルドステージ
FROM golang:1.24 AS builder

# 作業ディレクトリの設定
WORKDIR /backend

# Goモジュールをキャッシュ
COPY go.mod .
COPY go.sum .
RUN go mod download

# アプリケーションのソースをコピー
COPY . .

# アプリケーションをビルド
RUN go build -o main main.go

# 実行ステージ
FROM debian:stable-slim

# Cloud Run が 8080 ポートでリクエストを待つため 環境変数を設定
ENV PORT=8080
WORKDIR /app

# なぜか path を通さないとエラーになった
ENV PATH="/app:${PATH}"

# ビルド済みのバイナリをコピー
COPY --from=builder /backend/main /app/main

CMD ["./main"]
