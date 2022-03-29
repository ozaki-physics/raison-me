# crypto-assets
とりあえず イメージを起こしてみる  
## システム関連図
```mermaid
flowchart LR
  user(僕)
  line(LINE アプリ)

  user -- html, curl --- system
  line -- API --- system 
  user -- app --- line

  subgraph system [raison-me]
    direction TB
    service01(暗号資産の価格管理)
  end
  
  exchange(取引所)
  system -- API --- exchange
```
## ユースケース図
```mermaid
flowchart LR
  subgraph user [僕]
    direction LR
    browser(ブラウザ)
    console(コンソール)
    line(LINE アプリ)
  end

  subgraph system [暗号資産の価格管理]
    direction LR
    use01(平均約定レート を確認する)
    use02(評価額 を確認する)
    use03(保有枚数 を確認する)
    use04(コインの損益 を確認する)
    use05(損益の割合 を確認する)
    use06(価格ごとの枚数 を確認する)
    use07(ある取引の約定代金 を確認する)
    prefix01(1つのコイン)
    prefix02(すべてのコイン)
    prefix01 & prefix02 --- use01 & use02 & use03 & use04 & use05 & use06 & use07
  end

  exchange(取引所)
  db("永続化\n(DB, JSON)")

  user --- prefix01 & prefix02
  use04 & use05 --- exchange 
  system --- db

  memo01("1つを繰り返せばすべてになるが API リクエスト回数を少なくするため\nまとめて1リクエストにしたユースケースも作る")
  exchange -.- memo01
```
## ドメインモデル図
```mermaid
flowchart RL
  transaction -- 1, 0...n ---> coin

  subgraph aggregation01 [コイン集約]
    subgraph coin [コイン]
      direction LR
      coin_symbol(銘柄 : symbol)
      coin_bid("現在値(売値) : bid")
    end
  end

  subgraph aggregation02 [取引集約]
    direction LR
    subgraph transaction [取引]
      direction LR
      transaction_symbol(銘柄 : symbol)
      transaction_side(売買区分 : side)
      transaction_price_rate(約定レート : price_rate)
      transaction_size(約定数量 : size)
      transaction_fee(取引手数料: fee)
      transaction_time(約定日時 : time)
    end
  end

rule01("buy か sell")
rule01 -.- transaction_side
```
## オブジェクト図
```mermaid
flowchart RL
  transaction01 --> coin01
  transaction02 --> coin01

  subgraph coin01 [コイン01]
    direction LR
    coin01_symbol(銘柄: BTC)
    coin01_bid(現在値: 410万)

    subgraph coin01_memo [メモ]
      direction LR
      get_price_average("平均約定レート: (400*0.01+410*0.01)/0.02=405万")
      price_btc(評価額: 405万*0.02枚=8.1万)
      mai_all(保有数量: 0.01+0.01=0.02枚)
      gain("評価額損益: (410万-405万)*0.02枚=0.1万")
      gain_percent(評価額損益%: 0.1万/8.1万=1.23%)
    end
  end
  subgraph transaction01 [取引01]
    direction LR
    transaction01_symbol(銘柄: BTC)
    transaction01_side(売買区分: 買)
    transaction01_price_rate(約定レート: 400万)
    transaction01_size(約定数量: 0.01枚)
    transaction01_fee(取引手数料: 1)
    transaction01_time(約定日時: 2022/03/09 13:00)
    subgraph transaction01_memo [メモ]
      transaction01_price(約定代金: 400万*0.01枚=4万)
    end
  end
  subgraph transaction02 [取引02]
    direction LR
    transaction02_symbol(銘柄: BTC)
    transaction02_side(売買区分: 売)
    transaction02_price_rate(約定レート: 410万)
    transaction02_size(約定数量: 0.01枚)
    transaction02_fee(取引手数料: -3)
    transaction02_time(約定日時: 2022/03/10 00:36)
    subgraph transaction02_memo [メモ]
      transaction02_price(約定代金: 410万*0.01枚=4.1万)
    end
  end
```
## ER 図
```mermaid
erDiagram
  COIN {
    int id PK "raison-me の capital での ID"
    int CMC_id "CoinMarketCap の 銘柄 ID"
    string symbol_name "仮想通貨の名前(Bitcoin)"
    string create_role "作成した人, プログラムの識別子"
    timestamp create_time "作成した日時"
    string update_role "更新した人, プログラムの識別子"
    timestamp update_time "更新した日時"
  }
  NOW {
    int id PK "NOW テーブルで一意にするため"
    int coin_id "COIN テーブルとリレーションするため"
    long bid "現在値の価格"
    string create_role "作成した人, プログラムの識別子"
    timestamp create_time "作成した日時"
    string update_role "更新した人, プログラムの識別子"
    timestamp update_time "更新した日時"
  }
  TRANSACTION {
    int id PK "TRANSACTION テーブルで一意にするため"
    int coinId "COIN テーブルとリレーションするため"
    int side "売買区分 1:買, 2:売"
    long price "約定レート"
    int size "約定数量"
    int fee "取引手数料"
    timestamp time "取引した日時"
    string create_role "作成した人, プログラムの識別子"
    timestamp create_time "作成した日時"
    string update_role "更新した人, プログラムの識別子"
    timestamp update_time "更新した日時"
  }
  COIN ||--o{ NOW : "coin の id"
  COIN ||--o{ TRANSACTION : "coin の id"
```
## クラス図
mermaid から drawio で書いた [svg](./crypto-assets.svg) に変更  

ここには気になったことを書く  
- 命名がよくない  
-> cryptoAssetsPresen は ちゃんと apiController, htmlController などにした方がいい  
-> config は ちゃんと credential などにした方がいい  
- coinMarketCap の依存も外部なので 逆転させた方がいいかも?  
-> それを担うのが infra 層 だから 気にしなくていい?  
-> でも infra 層が分厚くなるというか 取得した Response を 加工する Logic がドメインから漏れている?  
-> infra 層 を難しく考えすぎてた笑 単純化できた  
- すると 現在値を coinMarketCap 以外から取得 や key を SecretManager 以外から取得 ができる?  
-> credential を保持する interface を用意して回避した  
- interface にするなら ドメイン層 まで持っていくべき?  
-> key value は infra 層でしか使わないから domain 層まで持っていかない  
- Presentation 層で 同じことをするって意味合いはどうやって表現するか?  
-> Usecase 層 の同じ interface を使っていることで表現できるかと  
- Usecase 層 は メソッド1個1個で interface で公開すべき?
-> 現時点ではあるまとまりで interface を作って 各 Controller に その interface を実装する
- Repository にDI するときに CMC か GMOCoin かを入れると切り替えができそう  
- LINE bot のために key が必要だが Secret Manager は infra 層 にあるから層を越えることになる  
  - SecretManager にアクセスするクラスは どの層からも使えるような場所を作る? -> 低凝集になる  
  - Presentation 層で SecretManager にアクセスするコードを書く? -> サービス内に同じことが書かれたコードが増え冗長  
  - credential interface を domain 層に持っていき Usecase 層を経由して使う? -> LINE bot で使えるかは使用方法の1種だからドメインほどコアではない  
-> key は環境変数みたいなものだからどこからアクセスしてもいい気がする, しかし コアに関係なく外界の環境のためにある  
-> よって infra 層の Share だけ どの層から呼んでも良いとする(絶対に各サービスを越える共通オブジェクトは作らない, 将来 各サービスごとに分離できなくなるから)  
- 中途半端に 本番環境と開発環境, テストとモック のことを考えているから冗長なコードが増えている?  
-> 勉強不足ではあるが とりあえず作りたいから作り進める  
- ファイル名の接頭辞に cryptoAssets_ が付与されていて cryptoAsset 以外のコンテキストを作ったときに名前空間が衝突する  
  - そもそもパッケージ名はアンダースコアや大文字が混ざるのは良くない  
  - 一般的に使われてる単語 crypto だけも良くない  
-> 妥協して cryptoAsset にする(domain/cryptoAsset など) 各層のファイルには `package cryptoAsset` と記述されるが諦める
-> 同じ `package cryptoAsset` と記述されるが ちゃんと名前空間は別になっている  
-> import するときに 別名で層の名前を付けて識別する  
-> 自分がどの層の cryptoAsset を編集しているかは ファイル名から判断する  
- capital/infrastructure/json にある LINE bot の key は LINE bot(テスト)  
-> LINE bot(本番)は GCP のSecret Manager  
-> テスト環境にデプロイしたら テスト環境 の SecretManager に保存された LINE bot(テスト)の key が使われるだけ  
- credential の interface は何を抽象化するのか  
-> 何の永続化ツール(db, JSON, GCP など)から入手した情報かを 区別しない  
-> ただ どのサービス(CoinMarketCap, LINE)のクレデンシャルかは 区別する  
- credential が テスト環境か 本番環境か 区別を抽象化したい  
-> 外部から infra 層で使いやすい形の Dto を作る  
-> その Dto を作るときに テスト環境か 本番環境か の boolen を渡して判別する  
- 正しい使い方じゃない気がするが 外部から取得して infra 層で使いやすい形に変換したものを XXXDto って名前つけがち  
- 詰め替えるときのメソッド名を generate にしがち(Create だと DI するときの init な雰囲気するから 使わない, make だと generate よりも軽く一部を抜き出しただけな感じ)  
- LINE の通信に必要な apiKey などは infra 層にする, ただ Presentation 層から infra 層へ依存するのはダメなので infra/Share という特別な package に配置する  
-> Presentation 層は 外部との内容に注力すべきであって 外部と通信するために必要なものは Presentation 層ではないと思うから infra 層にした  
- 命名規則というほどじゃないけど json をマッピングした struct には dataXXX って名前をつけがち  
- REST API にするために `http.HandleFunc()` の第2引数では http メソッドを固定することができない  
-> echo, gin は どうやって url の仕分けと http メソッド の仕分けを両立しているんだろう  
-> 勉強がしたいのか, 便利ツールを作りたいのか どちらを優先しようか
