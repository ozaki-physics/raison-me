# ADR 0003: 独自 JWT による認証方式を採用する
ステータス: Accepted

## 背景(コンテキスト)
`info/authN` には `signin` や `token/refresh` のルートが存在するが、認証トークンの設計方針はまだ固まっていない。
今後は `Authorization: Bearer <token>` で API 認証を統一したいが、ログインのたびに毎回再認証を要求すると UX が悪化しやすい。

一方で、長寿命トークンをそのまま使うと漏えい時のリスクが大きい。
そのため、短命な `access token` と、再発行専用の `refresh token` を分けるべきかを整理する必要がある。

歴史的には、従来の Web アプリはサーバが HTML を返し、ログイン状態は server-side session と cookie で維持する構成が中心だった。
この構成では、同一バックエンドの中で session ストアを参照できれば十分であり、認証情報を API ごとに自己完結した形で持ち運ぶ必要は小さかった。

しかし SPA、モバイルアプリ、フロントエンドとバックエンドの分離、外部 API 連携が一般化するにつれて、
画面表示中に API を連続して呼び出しつつ、再ログインなしで利用を継続したい要求が強くなった。
同時に、複数 API や将来のサービス分割を見据えると、各リクエストで共有 session ストアへ依存せずに認証済み情報を検証したい需要も増えた。

業界的にも OAuth 2.0 と OpenID Connect の普及により、短命な `access token` と長命な `refresh token` を分ける設計が広まった。
もともと OAuth 2.0 は認可、OpenID Connect は認証のための仕様だが、
このエコシステムを通じて「短命 access token で API を呼び、refresh token で再発行する」という運用パターンが一般的になり、独自認証でも同じ分離を採用するケースが増えた。

ただし `refresh token` をサーバ側で保存して失効管理し始めると、見た目としては session 管理に近づく。
それでも `access token` に JWT を使う意味が残るのか、同じ要件を opaque な session token で実現すると何が辛くなるのかを意思決定として記録しておく。

## 課題
以下を明確にする必要がある。

- `access token` と `refresh token` をどう役割分担するか
- `refresh token` を導入して stateful 要素が入る中で、JWT を使う意味がどこにあるか
- 同じことを server-side session / opaque token 方式で実現する場合、何が実装上・運用上の負担になるか
- 認証ミドルウェア、失効制御、将来の API 拡張をどの方針に寄せるか

## 決定
以下を採用する。

- `access token` は独自 JWT を使う
- `refresh token` は再発行専用の資格情報として扱い、サーバ側で保存、失効、ローテーション可能にする
- 保護された API は `Authorization: Bearer <access_token>` で認証する
- 認証ミドルウェアは `access token` の署名、期限、必須 claim を機械的に検証し、認証結果を context へ渡す責務に寄せる
- `signout` は `refresh token` または refresh session の失効で表現する
- `access token` の即時失効は基本要件にしない。漏えい時の影響は短い有効期限で抑える

## 記録
- 作成日: 2026-04-05
- 更新日: 2026-04-05
- 置き換え先 ADR: なし

## 検討した選択肢
1. `access token` のみを使い、期限切れ後は毎回再ログインさせる
2. `JWT access token` と、サーバ側で管理する `refresh token` を併用する
3. JWT を使わず、opaque な session token または server-side session に統一する

## 決定理由
JWT を使う主な理由は `refresh token` 側ではなく `access token` 側にある。
`access token` を JWT にすると、各 API は共有ストアへ毎回アクセスしなくても、署名検証だけで `sub`、`exp`、`role` などの認証済み情報を読み取れる。
これにより、認証ミドルウェアの責務を「Bearer token を取り出し、署名と期限を検証し、結果を context に詰める」という機械的な処理へ寄せやすい。

一方で `refresh token` は長寿命であるため、漏えい時の被害を抑えるには失効やローテーションが必要になる。
この部分を stateful に管理するのは、JWT の利点を否定するためではなく、長寿命資格情報だけを安全側へ寄せるための責務分離である。
短命な `access token` は自己完結性を優先し、長命な `refresh token` は管理可能性を優先する。

server-side session や opaque token 方式でも同じ UX は実現できるが、その場合は全リクエストで session ストアへの参照が必要になる。
すると認証ミドルウェアがストレージ依存になり、単なる検証器ではなく `session lookup`、期限判定、失効判定、必要に応じた user 取得まで抱えやすい。
この構成は単一バックエンドでは十分成立するが、次の点が辛くなりやすい。

- 認証可否が session ストアの可用性と応答時間に毎回依存する
- スケールアウト時に共有ストア設計が必須になる
- claim ベースの認可をしたいとき、session か user を追加で読み出す箇所が増えやすい
- 将来 API が分かれたとき、各サービスが session ストアへ結合しやすい
- ミドルウェアの責務が肥大化し、テストが「token 文字列の検証」ではなく「ストア状態込みの認証」に寄りやすい

今回のバックエンドは現時点では単一サービスであり、opaque session でも十分構築可能である。
それでも 以下を考えると、`JWT access token + stateful refresh token` が最もバランスが良い。
- `Authorization: Bearer` を軸に API 認証を揃えたいこと
- 今後 claim を使った `me` や認可判定を薄く書きたいこと
- 将来的なサービス追加余地を残したいこと

また、この選択は「最近の流行を追う」ためではなく、業界で増えた利用形態に対する自然な帰結でもある。
単一の server-side session で閉じる世界から、API を中心に複数のクライアントが継続利用する世界へ重心が移った結果、
短命 access token と refresh token を分離する方式が広く採用されるようになった。
本 ADR もその流れを踏まえつつ、現行バックエンドの規模と将来余地の両方に整合する設計として判断している。

## メリット: 期待される効果
- 短命な `access token` に自己完結性を持たせられる
- 認証ミドルウェアを署名検証中心の薄い責務にしやすい
- `Authorization: Bearer` 前提の API と相性が良い
- `refresh token` を失効、ローテーションできるため UX と安全性を両立しやすい
- 将来 API やサービスが増えても、認証済み情報の受け渡し方法を揃えやすい

## デメリット: 受け入れるリスク
- `refresh token` の保存、失効、ローテーションの実装が必要になる
- `access token` は基本的に即時失効しないため、短い有効期限設計が前提になる
- JWT の署名鍵管理と将来の鍵ローテーション方針が必要になる
- claim に情報を載せすぎるとトークン設計が硬直化する
- session 管理と比べると、token 種別ごとの責務分離を設計する手間が増える

## 影響範囲
- コード: `info/authN/usecase`、`info/authN/presen`、認証 middleware、config に `access token` / `refresh token` の責務追加が必要になる
- テスト: `signin`、`refresh`、`signout`、JWT 検証 middleware、refresh token の失効とローテーションに対するテストが必要になる
- 運用: JWT 署名鍵の管理、refresh token 保存先、期限ポリシー、失効監査の方針が必要になる

## 将来の考慮事項
- `refresh token` をランダム文字列にするか JWT にするか
- `refresh token` を保存時にハッシュ化するか
- JWT の claim をどこまで最小化するか
- 鍵ローテーション時に `kid` を導入するか
- 複数端末ログイン時に refresh session を端末単位で失効できるようにするか

## 参考
- RFC 7519: JSON Web Token (JWT)
- OWASP Cheat Sheet Series: Session Management Cheat Sheet
- OWASP Cheat Sheet Series: JSON Web Token for Java Cheat Sheet
