# サポート期間: 2025-08-xx まで
# https://go.dev/doc/devel/release
FROM golang:1.23 AS golang
# サポート期間: 2025-10-xx まで (最大 2027-04-30)
# https://nodejs.org/en/about/previous-releases
FROM node:22 AS nodejs

# ローカル開発用のイメージ作成
# サポート期間: 2029-04-xx まで
FROM ubuntu:24.04  AS local

# 開発で使いそうなツールをインストール
RUN apt-get update && \
    # それぞれの ライブラリ をインストールするとき
    # 推奨されるパッケージをインストールしない ときは --no-install-recommends オプションをつける
    # apt-get install -y --no-install-recommends
    apt-get install -y \
    build-essential \
    curl \
    git \
    && apt-get clean \
    # キャッシュ削除
    && rm -rf /var/lib/apt/lists/*

# Google CLI のインストール
# https://cloud.google.com/sdk/docs/install-sdk?hl=ja#deb
RUN echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | \
    tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && \
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | \
    gpg --dearmor -o /usr/share/keyrings/cloud.google.gpg && \
    apt-get update -y && \
    apt-get install google-cloud-sdk -y

# Go と Node.js から 動かすための必要な要素を コピー
COPY --from=golang /usr/local/go /usr/local/go
COPY --from=nodejs /usr/local/bin/node /usr/local/bin/node
COPY --from=nodejs /usr/local/lib/node_modules /usr/local/lib/node_modules

# Go の パスを通す
# Node.js は /usr/local/bin/ にコピーしたため すでにパスが通っている
ENV PATH="/usr/local/go/bin:${PATH}"

# npm を使うときに 生成させたら COPY だったが 動かなかった
# COPY --from=nodejs /usr/local/bin/npm /usr/local/bin/
# COPY --from=nodejs /usr/local/bin/npx /usr/local/bin/
# おそらく コピーしてきたときに パスに差異があるから
# なので コマンド本体 と bin 内のコマンド で エイリアス を作った
# node 公式イメージ内でも /usr/local/bin/npm -> ../lib/node_modules/npm/bin/npm-cli.js となっている
RUN ln -s /usr/local/lib/node_modules/npm/bin/npm-cli.js /usr/local/bin/npm
RUN ln -s /usr/local/lib/node_modules/npm/bin/npx-cli.js /usr/local/bin/npx

ENV REPOSITORY=app

# バックエンド の 設定
WORKDIR /$REPOSITORY/backend
# go.mod と go.sum の想定だが go.sum が無くても動くようにする
# .devcontainer ディレクトリの中だが compose.yaml で context を 更に上のディレクトリにしているため
# 相対パスで 更に上にいかなくて良い
COPY ./backend/go.* ./
RUN go mod download

# フロントエンド の 設定
WORKDIR /$REPOSITORY/frontend/raison-me-fe
COPY ./frontend/raison-me-fe/package*.json ./
RUN npm install

WORKDIR /$REPOSITORY

# 本当は ソースのコピー と 実行コマンド を定義しないといけない
# しかし この image は 開発環境 かつ docker compose からマウントして bash で 動かす前提のため
# ソースのコピー と 実行コマンド(CMD)も書かない
