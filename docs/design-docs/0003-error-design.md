# Design Document 0003: レイヤ間のエラー設計

## 対象
- service: backend
- feature: error-handling

## 目的
レイヤごとのエラー責務を明確にし、usecase 層の公開契約と、domain / infra からの失敗の扱いを統一する。

## 結論
- usecase 層の公開メソッドは `UsecaseError` を返す
- usecase 層が依存する port / interface は原則 `error` を返す
- domain 層の検証や生成で起きた失敗は `DomainError` を返してよい
- infra 層の実装は通常の `error` または infra 独自 error を返してよい
- usecase 層は下位層の error を受け取り、文脈を付けて `WrapUsecaseError(...)` で包んで返す
- `UsecaseError` を使う基準は「その関数で処理があるか」ではなく「usecase 層の外へ返す境界かどうか」で決める

## 記録
- 作成日: 2026-04-18
- 更新日: 2026-04-18
- 関連 ADR: なし

## スコープ
- 含む: backend の usecase / domain / infra 間のエラー受け渡し
- 含む: usecase 層の公開契約の決め方
- 含む: `info/authN/usecase/refresh_token.go` を例にした判断基準
- 含まない: HTTP status code や API response body の詳細設計
- 含まない: ログ集約基盤や監視基盤の設計

## 前提
- このリポジトリでは DDD を採用し、プレゼンテーション層は薄く保つ
- usecase 層は domain と infra を調停し、プレゼンテーション層へ DTO を返す
- エラーも同様に、下位層の失敗をそのまま露出するのではなく、usecase 境界で意味付けして返す

## ユースケース
- usecase 実装者が、新しい port の戻り値を `error` と `UsecaseError` のどちらにすべきか判断できる
- プレゼンテーション層実装者が、usecase 呼び出し時に `UsecaseError` を入口として扱える
- domain / infra 実装者が、usecase 層の型に依存せず実装できる

## 設計

### 全体像
- domain 層は業務ルールと値検証の失敗を表現する
- infra 層は永続化や外部依存の失敗を表現する
- usecase 層は下位層の失敗を受け取り、ユースケースの文脈に変換して公開する
- presentation 層は usecase 層が返した `UsecaseError` を扱い、下位層の詳細型へ直接依存しない

### レイヤごとのルール
- domain 層:
  - 値オブジェクト生成や業務ルール違反は `DomainError` を返す
  - usecase 層の `UsecaseError` には依存しない
- infra 層:
  - repository / gateway 実装は `error` を返す
  - 必要なら infra 独自 error を返してよいが、usecase 層の `UsecaseError` には依存しない
- usecase 層:
  - 公開メソッドは `UsecaseError` を返す
  - 下位層の `error` / `DomainError` を受け取り、`WrapUsecaseError(...)` で文脈を付ける
  - 依存先 port の戻り値として `UsecaseError` を要求しない
- presentation 層:
  - usecase から返る `UsecaseError` を元に response を構築する
  - domain / infra のエラー型判定を直接持ち込まない

### 判断基準
- `UsecaseError` を返すべき関数:
  - usecase 層の公開契約になる関数
  - usecase パッケージの外へ失敗を返す関数
- `error` を返すべき関数:
  - usecase 層が依存する port / interface
  - usecase 層の外で実装されることを前提にした抽象
- `DomainError` を返すべき関数:
  - domain の生成, 再構築, 検証ロジック

### authN の具体例
`info/authN/usecase/refresh_token.go` の `RefreshTokenGenerator` は次の理由で `error` を返すのが妥当とする。

- 実装側に usecase 層の型を背負わせないため
- usecase 層が infra / domain に依存する方向は許容されても、その逆方向の依存は避けるため
- 実装が返した失敗は、呼び出し側 usecase で `WrapUsecaseError(...)` すれば十分なため

同ファイルの `NewRefreshToken(...)` は usecase パッケージの外へ返す補助関数として扱うなら `UsecaseError` を返してよい。  
ただし理由は「処理があるから」ではなく、「usecase 層の公開契約として返すエラーを統一するため」と説明する。

## エラー設計
- 下位層の失敗を usecase 層でそのまま返さない
- usecase 層では、呼び出し元が判断しやすい文脈をメッセージに含める
- invalid input や not found のように公開上まとめたい失敗は、sentinel error を inner error にして `WrapUsecaseError(...)` する
- ログやデバッグで元原因が必要なため、`Unwrap()` で内側の error を辿れる設計を維持する

## アンチパターン
- 依存先 interface の戻り値を `UsecaseError` にする
- 「エラーになる可能性があるから」という理由だけで usecase 専用 error を選ぶ
- domain / infra の詳細 error を presentation までそのまま露出する

## テスト観点
- usecase 層の公開メソッドが `UsecaseError` を返す
- domain の生成失敗を usecase 層で `WrapUsecaseError(...)` できる
- infra の失敗を usecase 層で `WrapUsecaseError(...)` できる
- `errors.Is` / `errors.As` で inner error を辿れる
- presentation 層が `UsecaseError` だけを入口に扱える

## 影響範囲
- コード: `info/*/usecase`, `info/*/domain`, `info/*/infra` のエラー設計方針
- テスト: usecase error の wrap / unwrap を確認するテスト
- 運用: 新しいレイヤ追加時の設計判断基準

## 未解決事項
- `UsecaseError` と `DomainError` のメッセージをどこまで日本語 / 英語で統一するか
- API 向けのエラーコード設計を別ドキュメントとして切り出すか

## 参考
- [0002-authn-token](./0002-authn-token.md)
- [refresh_token.go](../../backend/info/authN/usecase/refresh_token.go)
