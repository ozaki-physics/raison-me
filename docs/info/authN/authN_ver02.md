# 認証管理
[よくわかる認証と認可](https://dev.classmethod.jp/articles/authentication-and-authorization/)  
認証(Authentication: AuthN)  
認可(Authorization: AuthZ)  

アクセストークン は 認可(authZ)のコンテキストとする  

一度全部作り直す  
なぜなら 認証をパスワード以外も対応したいから  
## ドメインモデル図
```mermaid
flowchart LR
  subgraph user [ユーザー]
    direction LR
    accountID(アカウントID)
    userID(ユーザーID)
    userName(ユーザー名)
  end

  subgraph pass [パスワード]
    direction LR
    passID(パスワードID)
    password(パスワード)
    pass_iat(発行日時)
  end

  subgraph api_key [APIキー]
    direction LR
    keyID(キーID)
    key(キー)
    key_iat(発行日時)
  end

  subgraph profile [プロファイル]
    direction LR
    profileID(プロファイルID)
    description(備考)
  end

profile -- 1 ... 1 --> user
password -- 1,n ... 1 --> user
api_key -- 0,n ... 1 --> user
```

プロファイル は過剰機能だから 最初は実装しない  
## 開発中のメモ
ドメイン層 で発生する例外を作らないとなぁ  
enum の 作り方とかあったなぁ  
テストコードはインタフェースに対して作るのがよい  
ログ出力ミドルウェア  
コンフィグファイル  

値オブジェクトのファイル名を xxx_value.go にしたのは失敗だったかも  
テストファイルが xxx_test.go って書くから めちゃめちゃ長くなる  
