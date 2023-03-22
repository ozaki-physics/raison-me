# 使い方 Usage
## Remote - Containers
VS Code の拡張機能 Remote - Containers(識別子: ms-vscode-remote.remote-containers) を使って開発する  
コンテナ内で VS Code を起動し go 言語のための VS Code の拡張機能 Go(識別子: golang.go) を使う  
[golang.go](https://marketplace.visualstudio.com/items?itemName=golang.Go)  

golang.go は 様々なモジュールをインストールする必要がある  
特に gopls が go のコード補完してくれる便利なライブラリで使いたい  
ただ そのモジュールを Docker image に含めないようにするため  
コンテナを起動した後の コンテナ内 VS Code で変更を加える  
.devcontainer/devcontainer.json で gopls だけ install しているが VSCode でコンテナに入った時は認識してくれないから  
install が終わったら gopls だけ VSCode の通知から install にすればよい  

コンテナ内の git では日本語が使えないため コミットするときは ローカルの git bash 等を使う  
```bash
$ docker-compose build
$ docker-compose up -d
# VS Code より Remote - Containers で接続する
# VS Code の通知(golang.go)より gopls だけ install をする

# 基本は VS Code 内のターミナルで良いが ローカルの PowerShell からアクセスしたくなった場合
$ docker-compose exec raison_me bash
# 初めて開発を始めるときは go mod を生成する
/go/src/github.com/ozaki-physics/raison-me# go mod init $REPOSITORY
# 終えるとき
# VS Code より Remote - Containers で接続をやめる
$ docker-compose down
```

ちなみに image の時点で go build は済んでおり  
image から直接 run または docker-compose.yml の command をコメントアウトで確認できる  
```bash
$ docker container run --rm -d -p 8080:8080 --name go_raison_me go1.16:raison_me
$ docker container stop go_raison_me
```

### 外部モジュールのバージョンアップ
例として github.com/gin-gonic/gin をバージョンアップする  
1. コンテナにアタッチする
2. モジュールのバージョンアップして go.mod を更新する
3. コンテナを削除してもバージョンアップが反映されるように docker image を作り直す

```bash
$ docker-compose up -d
$ docker-compose exec raison_me bash

# モジュールのバージョンアップ
/go/src/github.com/ozaki-physics/raison-me# go install -v github.com/gin-gonic/gin
# 不要モジュールの削除
/go/src/github.com/ozaki-physics/raison-me# go mod tidy -v

$ docker-compose down
$ docker image rm go1.16:raison_me
# docker image の作り直し
$ docker-compose build
```

### Go のバージョンを上げる
1. Dockerfile  
Go の バージョンを上げる  
`RUN go mod download` をコメントアウト  
2. docker-compose.yml の image 名を変更  
3. image を build  
4. コンテナ起動して exec で接続  
5. go.mod の中身をほぼ空にする  
以下だけにする  
```
module github.com/ozaki-physics/raison-me

go 1.19
```
6. `go mod tidy` を実行  
すると go.mod と go.sum が変更される  
7. app.yaml の GAE のバージョンを上げる  
8. Dockerfile で `RUN go mod download` のコメントアウトを戻す  
9. image を作り直す  
`docker image rm go1.19:raison_me`  
`docker-compose build`  

## gcloud docker
[GCP Container Repository](https://console.cloud.google.com/gcr/images/google.com:cloudsdktool/GLOBAL/cloud-sdk?authuser=9)をもとに pull する image を決める  
```bash
$ docker pull gcr.io/google.com/cloudsdktool/cloud-sdk:423.0.0
# カレントディレクトリの内容を全部マウントする
$ docker run -it --rm -v "$(pwd):/raison-me" gcr.io/google.com/cloudsdktool/cloud-sdk:423.0.0 bash
# 認証する
$ gcloud auth login --no-browser
# 表示される コマンドを gcloud がインストールされたホストのターミナルで実行
# ブラウザが起動するからアカウントを許可(認証が完了すると勝手に遷移する)
# ホストのターミナルで表示された URL を コピーして コンテナ内のターミナルに入力
# 認証される
$ cd raison-me/
# プロジェクト ID の設定
$ gcloud config set project PROJECT_ID
# Google Compute Engine 使用されるデフォルトのゾーンを設定
$ gcloud config set compute/zone us-west1-a
# デフォルトのリージョンを設定
$ gcloud config set compute/region us-west1
# GAE をデプロイ
$ gcloud app deploy
```
または  
```bash
$ gcloud init
# GCP のプロジェクトを選ぶ(raison-me を選択)
$ gcloud app deploy
# リージョンは us-west1
```

## SecretManager
[Secret Manager client libraries](https://cloud.google.com/secret-manager/docs/reference/libraries#client-libraries-install-go)  
`$ go get cloud.google.com/go/secretmanager/apiv1`  
`go get: added cloud.google.com/go/secretmanager v1.3.0`  
`$ go get google.golang.org/genproto/googleapis/cloud/secretmanager/v1`  
`$ go mod tidy`  
