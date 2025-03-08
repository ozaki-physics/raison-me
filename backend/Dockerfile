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

RUN apt-get update && \
    apt-get install -y \
    ca-certificates

WORKDIR /app

# なぜか path を通さないとエラーになった
ENV PATH="/app:${PATH}"

# ビルド済みのバイナリをコピー
COPY --from=builder /backend/main /app/main

# persist(DB の代わり) に使う データ を コピー
COPY --from=builder /backend/capital/infrastructure/cryptoAsset/json/coinMarketCapId.json /app/persist/coinMarketCapId.json
COPY --from=builder /backend/capital/infrastructure/cryptoAsset/json/transaction.json /app/persist/transaction.json

CMD ["./main"]
