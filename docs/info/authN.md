# アカウント管理
アカウント管理 じゃなくて 認証管理?

IDトークン に 有効無効 区分は必要?
レコード削除はいささか乱暴? 論理削除にする? データ量削減しちゃう?
有効期限を付ける?
どこで使うために発行された IDトークン か分かるようにする?
備考を付ける?
JWT の規格に合わせる?
[よくわかる認証と認可](https://dev.classmethod.jp/articles/authentication-and-authorization/)
認証(Authentication: AuthN)
認可(Authorization: AuthZ)
## システム関連図
## ユースケース図
- ユーザー を新規作成できる
- ユーザー を削除できる
- ユーザー でログインできる
- ユーザー(自分) の情報を更新できる
  - ユーザーID
  - ユーザー名
  - パスワード
  - IDトークン
  - アクセストークン
- IDトークン の 発行 ができる
- IDトークン の 無効化 ができる
- アクセストークン の 発行 ができる
- アクセストークン の 無効化 ができる

まずは __限られた人にしか見れない__ という状態にする  
## ドメインモデル図
```mermaid
flowchart LR
  subgraph user [ユーザー]
    direction LR
    accountID(アカウントID)
    userID(ユーザーID)
    userName(ユーザー名)
    password(パスワード)
  end
  subgraph token [トークン]
    direction LR
    tokenID(トークンID)
    IDToken(IDトークン)
    accessToken(アクセストークン)
  end

token -- 0,1 ... 1 --> user

rule01("日本語, 絵文字 OK")
rule01 -.- userName
rule02("文字種の制限")
rule02 -.- password
rule03("表には出さない")
rule04("管理者を含めても一意")
rule03 & rule04 -.- accountID
```

ユーザーID を変えられるように アカウントID を作った  
だけど ユーザーID も一意にならないと メンションができない  
そもそも一意になりそうな値をハッシュ化する?  
アカウントID は表に出さない値だから どんなのでもいい  

ついつい DB のことまで考えて key とか どっちのテーブルにIDを持つか とか考えてしまう  
あくまで オブジェクト として考えないといけない  
## オブジェクト図
オブジェクト図は 具体的値, ドメインモデル図 は英語名 も書く  
出てくる名詞をとりあえず オブジェクト として考えてみるといい感じ  
## ER 図
```mermaid
erDiagram
  ACCOUNT {
    int id PK "ACCOUNT テーブルで一意にするため"
    string user_id "ユーザーID"
    string user_name "ユーザー名"
    string password "パスワード"
    string create_role "作成した人, プログラムの識別子"
    timestamp create_time "作成した日時"
    string update_role "更新した人, プログラムの識別子"
    timestamp update_time "更新した日時"
  }
  TOKEN {
    int id PK "TOKEN テーブルで一意にするため"
    stirng account_id "ACCOUNT テーブル とリレーションするため"
    string id_token "IDトークン"
    string access_token "アクセストークン"
    string create_role "作成した人, プログラムの識別子"
    timestamp create_time "作成した日時"
    string update_role "更新した人, プログラムの識別子"
    timestamp update_time "更新した日時"
  }

  ACCOUNT ||--o| TOKEN : "ACCOUNT の id"
```
## クラス図
