# Design Document 0001: AuthN Token

## 対象
- service: info
- feature: auth-n

## 目的
`docs/adr/0003-custom-jwt-authentication-strategy.md` の方針に沿って、最小限の DB 設計で `signin`、`token/refresh`、`signout` を実現できるようにする。

## 結論
- `access token` は DB に保存しない
- 既存の `app.accounts` と `app.passwords` を使う
- 新しく `app.refresh_tokens` を追加し、refresh token だけ DB で管理する
- refresh token は hash で保存する
- 1 account で複数 refresh token を許可し、複数端末利用に対応する
- refresh 時は対象 token の row だけを更新して新しい token に差し替える
- signout 時は `revoked_at` を埋めて失効させる

## 記録
- 作成日: 2026-04-05
- 更新日: 2026-04-05
- 関連 ADR: ADR 0003 [custom-jwt-authentication-strategy](../adr/0003-custom-jwt-authentication-strategy.md)

## スコープ
- 含む: refresh token を保存する最小限のテーブル設計
- 含む: `signin`、`refresh`、`signout` で必要になる基本カラム
- 含む: 1 account で複数 refresh token を持てる設計
- 含まない: access token blacklist
- 含まない: 端末ごとの詳細監査
- 含まない: token reuse 検知などの高度な制御

## 前提
- `access token` は JWT であり、API 認証時に DB lookup しない
- `refresh token` は長寿命なので DB で失効管理する
- `refresh token` の平文は DB に保存しない
- 既存の認証主体は `app.accounts.account_id` を使う
- 既存のパスワード検証は `app.passwords` を使う

## ユースケース
- ユーザーが `signin` すると access token と refresh token を受け取れる
- ユーザーが `refresh` すると新しい access token と refresh token を受け取れる
- ユーザーが複数端末で同時にログインできる
- ユーザーが `signout` すると、その refresh token だけが使えなくなる

## 設計

### 全体像
- `accounts`: 認証主体
- `passwords`: サインイン時の本人確認
- `refresh_tokens`: refresh token の保存と失効管理

保護 API は JWT の検証だけを行う。  
DB を参照するのは `signin`、`refresh`、`signout` のときだけにする。  
`refresh_tokens` は account 単位ではなく token 単位で持ち、1 account に複数 row を許可する。

### 処理フロー
1. `signin` で `accounts` と `passwords` を使って本人確認する
2. access token と refresh token を発行する
3. refresh token を hash 化して `refresh_tokens` に保存する
4. `refresh` では token を hash 化して照合し、対象 row の token を更新する
5. `signout` では対象 row の `revoked_at` を更新して失効させる

## DB
### 概要
- `app.accounts`: 既存のまま利用
- `app.passwords`: 既存のまま利用
- `app.refresh_tokens`: 新規追加

### カラム / フィールド設計
| name | type | required | default | note |
| --- | --- | --- | --- | --- |
| refresh_token_id | UUID | yes | なし | 主キー |
| account_id | UUID | yes | なし | `accounts.account_id` への外部キー |
| token_hash | TEXT | yes | なし | refresh token のハッシュ値 |
| expires_at | TIMESTAMPTZ | yes | なし | 有効期限 |
| revoked_at | TIMESTAMPTZ | no | NULL | 失効日時 |
| created_at | TIMESTAMPTZ | yes | なし | 作成日時 |
| updated_at | TIMESTAMPTZ | yes | なし | 更新日時 |

### インデックス / 制約
- `refresh_token_id`: primary key
- `account_id`: foreign key -> `app.accounts(account_id)`
- `token_hash`: unique
- `account_id` index

### ライフサイクル
- 作成: `signin` 成功時に 1端末 1row として insert する
- 更新: `refresh` 成功時に対象 row の `token_hash`, `expires_at`, `updated_at` を更新する
- 更新: `signout` 時に対象 row の `revoked_at`, `updated_at` を更新する
- 削除: 期限切れデータは後で cleanup する

### 整合性
- `refresh` は 1 transaction で token の検証と更新を行う
- `token_hash` を unique にして同じ token の重複保存を防ぐ
- `account_id` は unique にしない。これにより 1 account で複数 token を持てる
- `revoked_at IS NULL` かつ `expires_at > NOW()` の row だけを有効とみなす

## バックエンド
### プレゼンテーション層
- request: `signin` は `userID`, `password`
- request: `refresh` は refresh token
- request: `signout` は refresh token
- param: userID, password, refresh token
- response: access token, refresh token
- entrypoint: `/info/auth-n/v1/signin`, `/info/auth-n/v1/token/refresh`, `/info/auth-n/v1/signout`

### ユースケース層
- usecase: `SignIn`
- input param: `userID`, `password`
- output dto: access token, refresh token
- role: 本人確認し token を発行して refresh token の hash を保存する

- usecase: `RefreshToken`
- input param: refresh token
- output dto: access token, refresh token
- role: refresh token を hash で照合し、対象 token を更新する

- usecase: `SignOut`
- input param: refresh token
- output dto: なし
- role: 指定された refresh token だけを失効させる

### ドメイン層
- entity / value object / domain service: `User`, `Pass` に加えて `RefreshToken` を追加する
- rule: 失効済みまたは期限切れの refresh token は使えない
- rule: 1 account は複数の refresh token を持てる

### インフラ層
- repository / gateway: `RefreshTokenRepo`
- 使用技術: PostgreSQL, pgx

### DTO
- `SignInResultDto`: access token, refresh token
- `RefreshTokenResultDto`: access token, refresh token

## フロントエンド
### 画面 / コンポーネント構成
- page: 該当なし
- component: 該当なし
- role: フロント実装はこの設計の対象外

### 画面遷移 / 状態遷移
- 該当なし

### 入出力設計
- input: 該当なし
- output: 該当なし
- validation: 該当なし

### 権限制御
- 該当なし

## エラー設計
- `signin` で user が見つからない: ログイン失敗として扱う
- `signin` で password 不一致: ログイン失敗として扱う
- `refresh` で token が見つからない: 再発行失敗として扱う
- `refresh` で revoked / expired: 再発行失敗として扱う
- `signout` で token が見つからない: 正常終了扱いでもよいが実装時に決める

## テスト観点
- `signin` 成功時に `refresh_tokens` へ hash が insert される
- 同じ account で複数回 `signin` すると複数 row が作られる
- `refresh` 成功時に対象 row の `token_hash` と `expires_at` が更新される
- `refresh` で revoked token は失敗する
- `refresh` で expired token は失敗する
- `signout` で対象 row の `revoked_at` が更新され、他の row には影響しない
- プレゼン層 と DTO / param の受け渡し
- 運用: 期限切れ row の cleanup 方法は後続で決める

## 影響範囲
- コード: `info/authN/usecase`, `info/authN/presen`, `info/authN/infra`, `info/authN/domain`
- テスト: `signin`, `refresh`, `signout`, repository の追加テスト
- 運用: refresh token の hash 方式、期限ポリシー、cleanup 方針の追加

## 未解決事項
- refresh のたびに同じ row を更新するか、新しい row を追加するか
- 端末単位の signout をどこまで厳密にやるか
- 端末を識別する `device_id` や表示用の端末名を保持するか
- ユーザーが自分の端末一覧を見て、端末ごとに refresh token を無効化できる機能を入れるか
- reuse 検知や session family を入れるか
- `user_agent`, `ip_address`, `client_name` などの監査情報を持つか
- access token に `jti` を入れるか

## 参考
- [ADR 0003](../adr/0003-custom-jwt-authentication-strategy.md)
- [db-connection](./db-connection.md)
