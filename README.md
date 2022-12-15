# raison-me
自分でアプリを作ってみる  

## サービス
### capital
- [暗号資産の価格管理](./docs/capital/crypto-assets.md)  
  - [クラス図みたいなもの](./docs/capital/crypto-assets.svg)  
### delight
### growth
### info
- [権限管理](./docs/info/auth-control.md)  
  - [オブジェクト図みたいなもの](./docs/info/auth-control_object-graph.svg)  
  - [ドメインモデル図みたいなもの](./docs/info/auth-control_domain-graph.svg)  
- [アカウント管理](./docs/info/account-manager.md)  
### regung
### seed
### zeit

## 環境 Requirement
- Docker

## 使い方 Usage
[詳細](./docs/usage.md)
1. VS Code の拡張機能 Remote - Containers で開発する  
gopls だけ VSCode の通知から install する  
2. 開発が終わったら Remote - Containers を閉じる  
3. gcloud が使えるコンテナを開発したコードをマウントしながら起動して アタッチする  
4. 開発したコードをデプロイする  
