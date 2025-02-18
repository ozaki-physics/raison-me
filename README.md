# raison-me
自分でアプリを作ってみる  

## サービス
### capital
- [暗号資産の価格管理](./tmp_stash/docs/capital/crypto-assets.md)  
  - [クラス図みたいなもの](./tmp_stash/docs/capital/crypto-assets.svg)  
### delight
### growth
### info
- [認証管理](./tmp_stash/docs/info/authN.md)  
- [認可管理](./tmp_stash/docs/info/authZ.md)  
### regung
### seed
### zeit

## 開発環境 Requirement
- Docker  

## 本番環境 Production environment
- Google Cloud  
<!-- TODO: 使ってるサービス1個ずつ表記するとか SVG で作りたい -->

## 使い方 Usage
[詳細](./tmp_stash/docs/usage.md)
1. VS Code の拡張機能 Remote - Containers で開発する  
gopls だけ VSCode の通知から install する  
2. 開発が終わったら Remote - Containers を閉じる  
3. gcloud が使えるコンテナを開発したコードをマウントしながら起動して アタッチする  
4. 開発したコードをデプロイする  
