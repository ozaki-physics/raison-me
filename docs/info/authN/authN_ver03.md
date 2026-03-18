# 認証管理  
[よくわかる認証と認可](https://dev.classmethod.jp/articles/authentication-and-authorization/)  
認証(Authentication: AuthN)  
認可(Authorization: AuthZ)  

アクセストークン は 認可(authZ)のコンテキストとする  

認証の設計を考え直す  
ver01 は とりあえず 考えてみた  
ver02 は パスワード以外の 認証方法 もできるように 拡張性を持たせる  
ver03 は DB の仕組みが導入された  
  
## ドメインモデル  
### ユビキタス言語  
- アカウント: ユーザー と 管理者 の 総称  
- ユーザー:   
- 管理者:   
- ID トークン:   
  
## ユースケース("誰" が "何" を "行動" できる)  
- ユーザー が アカウント を 新規作成 できる  
- ユーザー が アカウント を 削除 できる  
- ユーザー が アカウント を ログイン できる  
- ユーザー が アカウント を ログアウト できる  
- ユーザー が 自分のアカウント情報 を 更新 できる  
- ユーザー が ID トークン を 発行 できる  
- ユーザー が ID トークン を 無効化 できる  
  
## ER 図  
[ER 図](./authN_ver03.drawio.svg)  
  
## 各層の責務の定義が曖昧だな  
プレゼン層は 外界との接点, ユースケース層が Web でも console でも 気にせず使えるようにするための緩衝材  
ドメインオブジェクトの生成は なるべく ユースケース層 に寄せる  
  
  
## 開発中のメモ  
### アカウント ID と ユーザー ID がある理由  
ユーザー ID を 変えられるように アカウント ID を作った  
イメージは X (Twitter)  
ユーザー ID は メンションのために 一意 にする  
アカウント ID は ユーザーには表示しないから どんな値でもいい  
### アカウント ID は UUID にするか?  
容量節約なら 連番, 連番でよいという人も多い  
NewSQL 的には 十分にランダム(UUID v4)の方がホットスポットが生まれない  
UUID v7 も気になる  
-> 独自の ID の形を作った "文字-文字"  
なんで そんなことした?  
UUID だと長いし読みにくい, 数字だと システム全体の一意性が保てない  
全部コントロールできる方が 面白そうだったから笑  
最悪 苦しむことになっても そういう経験があるって言えるかな笑  
### IDトークン に 有効無効区分 は必要か?  
IDトークン を使わなくなったとき 論理削除 or 物理削除 どちらにするか?  
-> 物理削除 だと ストレージを減らせるし 何度も作り直せるから物理削除にしよう  
### トークンオブジェクト に 有効期限, 誰が発行したか, 備考 を保存できるようにするか?  
JWT の規格を確認してみよう  
[IDトークンが分かれば OpenID Connect が分かる](https://qiita.com/TakahikoKawasaki/items/8f0e422c7edd2d220e06)  
### ドメインオブジェクト と DB 設計  
ついつい DB のことまで考えて key とか どっちのテーブルにIDを持つか とか考えてしまう  
あくまで オブジェクト として考えないといけない  
### プロファイル で分けた理由は?  
ユーザーテーブルは ユーザーに関する項目だからって ついつい過剰になりやすいから  
### 未整理  
- DB の テーブル名 は 複数形 にすべき?  
- ドメイン層 で発生する例外を作らないとなぁ  
- enum の 作り方とかあったなぁ  
- テストコードはインタフェースに対して作るのがよい  
- ログ出力ミドルウェア  
- 値オブジェクトのファイル名を xxx_value.go にしたのは失敗だったかも  
テストファイルが xxx_test.go って書くから めちゃめちゃ長くなる  
- オブジェクト図は 具体的値, ドメインモデル図 は英語名 も書く  
出てくる名詞をとりあえず オブジェクト として考えてみるといい感じ  

- 費用対効果を考えて 何を 値オブジェクト にするか?  
  - Account の ID  
  - userID  
  - userName  
  - password  
- entity の 戻り値は ポインタがいいのか?  
- 引数は同じ型なら省略して書けるけど 後で型を変える可能性があるから分けておこうかな  
- struct を宣言するとき フィールド名 を省略しないほうがよい 省略してて フィールド増えるとエラーになる  
  -> むしろエラーになった方がいい? ゼロ値を設定するからエラーにならなくてもいい?  
- テストコード どこにどうやって書こう?  
- date_value をテキトーに string にしてるから ちゃんと go の時間を扱うオブジェクトにしたい  
- 値オブジェクト の生成責務は entity の コンストラクタ or usecase の利用側 の どちらだろう  
  -> それが ドメイン知識 ならドメイン層で プログラミング上なら どちらでもいいのでは  
- 戻り値が ポインタじゃないなら nil を返せず ゼロ値 を返さなければならない  
  - err を返すときに nil を返せないのは不便?  
  -> とりあえず Repository の 戻り値を ポインタ にして nil を返せるようにする  
- 独自例外も定義したい  
- アカウントとか nil, err を返していいの?  
- やっぱり user オブジェクトにして アカウントID を保持させつつ それを key にするって方がドメインに合いそう  
- やっぱり userID を key にした方が良さそう, でも ユーザーID が変更できなくなる...  
- ユーザー に 最終ログイン日 を持たせたい  
- セキュリティの観点の話始める?  
- アクセス IPアドレス とか保持した方がいい?  
- ユーザーID はシステム上の制限じゃなくて 機能的意味合い  
- Repository に 保存する って作ると Delete と少し意味が被る? json で操作してるからか?  
- infra 層でポインタを返すのは危険では? だって 同じ構造体を触る可能性がある  
- Repository で interface の合成を使うと 何系の処理なのか分かりやすくなる気がする  
- 合成もとは private な interface でいいかも  
- 値オブジェクト を作りすぎた  
  - Go を使っているのだから 別ファイルにせず 関連性の強いファイル内で書けばよい  
  - 値オブジェクト は 非可換 なものを扱うから 属性値 は向いてないかも  
  -> 値オブジェクト にしたいのか 型安全 が欲しいだけか ちゃんと見極めること  
- interface を細かく作りすぎた  
  - interface の合成は そんなに使わない方向で考えるほうが良さそう  
- ユーザー と ユーザープロファイル と 認証情報 は分けた方がいいらしい  
- json を読み書きする方法は複数種類あるけど どれがいいだろう?  
- 迷っている chi を読んで車輪の再発明をするか chi を使ってしまうか  
  -> 自分へのサービス提供を優先するなら chi を使ってしまった方がいい...  
  -> chi を読んでみたって記事を書くなら 読んでもいいかもしれんけど まだ記事を書くつもり無いから 今は chi を使うか笑  
  - やってることは ServeMux 作って Handler 実体は `ServeHTTP(ResponseWriter, *Request)` を登録して 呼び出すだけ  
  - フツーは pattern と `http.HandleFunc` を登録するのを 3項目目 として http メソッド を足すだけ  
  -> 標準ライブラリだけで 作ろうと 色々考えたけど やっぱり面倒くさくなったから chi を使う笑 __でもいつか chi を読んで記事を書きたい__  
- chi は 異なるパッケージで同じライブラリじゃないと動かないらしい 片方が "github.com/go-chi/chi/v5" もう片方が "github.com/go-chi/chi" なら 認識されない  
  - セマンティックバージョニング だから v4 と v5 は後方互換 が無く 同じ名前の別パッケージ  
  -> よって import でも v5 って明記した方が下手にバグるより前に 修正が必要とエラーになる  
- chi の導入の判断が早すぎた? 必要に迫られてからでもよかった? だって 標準ライブラリ から chi にするのは苦労しない  
- 標準ライブラリ で サブルーター を作るのはどうするんだろう  
- 1URL = 1構造体 = 1 ファイル = http メソッド分の usecase でいいのか?  
- chi に合う 認証のミドルウェア 作らないとな  
- chi で Post で値を取り出す方法は?  
- Request の値を entity に変換するのは どの層の役割か? presen 層? usecase 層?  
  -> バリデーション はどの層で? じゃなくて 層の責務に応じたバリデーションをする  
- Go 言語の flag パッケージ を上手に使えないかな?  
- 疑似 PUT にする リクエストヘッダー の実装(Java のため)  
- URL に対して JSON 構造を送ってくるときは body = JSON か!w  
  - body = JSON なら `Content-Type: application/json` にしないといけない  
- なんか salt するための値を取ってくる方法が上手じゃない気がするけど とりあえず 作る笑  
  __UserName で salt しちゃダメじゃね?w__  
#### なんか Post して Go で受けとるのが上手にできなかった  
POST のときに 以下で試したけど `application/x-www-form-urlencoded` じゃないとうまくいかないっぽい  
- Content-Type: application/json
- Content-Type: multipart/form-data
- Content-Type: application/x-www-form-urlencoded

[Content-Typeの一覧](https://qiita.com/AkihiroTakamura/items/b93fbe511465f52bffaa)  

なんか色々受け取り方があるらしい?  
```go
r.Form.Get("category")
r.FormValue("category")
r.PostFormValue("category")
```
どれも `r.ParseForm()` を実行した後じゃないとダメ?  

ParseForm をしないで post の値は取れない?  
postする側(htmlまたはjs)が multipart/form-data だと値は取れない?  
multipart/form-data の場合 ParseMultipartForm で処理する?  

URL パラメータだと どれでも値が取れる  
渡してるデータが json の形してるのも問題?  

[フォーム](https://www.twihike.dev/docs/golang-web/forms)  
- Form  
- PostForm  
- MultipartForm  
- FormValue  
- PostFormValue  
- FormFile  
`r.Form()` と `r.PostForm()` には `r.ParseForm()` が必要  
`r.MultipartForm()` には `r.ParseMultipartForm()` が必要  
通称フォーム は `Content-Type: application/x-www-form-urlencoded` らしい  
`r.Form()` と `r.PostForm()`
- `application/x-www-form-urlencoded` しか対応してない  
- key の列挙ができる  

`r.Form()` は URLパラメータ と name属性の値 があったら 両方送信されてくる  
`r.PostForm()` は name属性の値 で上書きされる  

`r.FormValue()` と `r.PostFormValue()`
- `r.ParseForm()` や `r.ParseMultipartForm()` を自動でやってくれる  
- `application/x-www-form-urlencoded` と `multipart/form-data` に自動対応  
- key の列挙ができない  
- 同じ key に複数の値があっても 最初の値しか取得してくれない  
- FormValue は フォーム で で name属性の値 が URLクエリ文字列 より 優先  
- FormValue は マルチパートフォーム で URLクエリ文字列 が name属性の値 より 優先  
- PostFormValue は name属性の値 が優先  

#### struct に使う単語のニュアンス整理  
定義: この struct の フィールド名 は〇〇で 型は〇〇で... って感じ  
宣言: メソッド内で struct に値を入れて変数にする感じ?  
生成: New な感じ, コンストラクタ に近い  
格納: setter な感じ  
### Go の ハッシュアルゴリズム
`bcrypt.GenerateFromPassword(bPass, bcrypt.DefaultCost)` は 自動で ソルトを付与してくれるらしい  
bcrypt(OpenBSD の EksBlowfish ベースのパスワードハッシュ)  
sha256 は 高速ゆえに パスワード保護には不向き  
今は argon2 がよいらしい  
Argon2 には i / d / id がある
id は両方の性質を混ぜてバランス型  
