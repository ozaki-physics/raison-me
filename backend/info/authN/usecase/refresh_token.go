// RefreshTokenGenerator と RefreshTokenHasher を分ける理由は、生成とハッシュ化の責務を分離するためです。
// 生成は新しいトークンを作成することに専念し、ハッシュ化はトークンのセキュリティを確保することに専念します。
// これにより、コードの可読性と保守性が向上し、将来的に異なるハッシュ化アルゴリズムを導入する際にも柔軟に対応できます。
package usecase

import (
	"time"

	"github.com/ozaki-physics/raison-me/info/authN/domain"
)

// リフレッシュトークンの情報
type RefreshToken struct {
	token     string
	expiresAt domain.Date
}

// リフレッシュトークンの発行
// (New メソッドの代わり)
// interface にしている理由は リフレッシュトークン の 発行方法 を 柔軟にしておくため
// 発行方法 は 技術的な概念 になるので infra 層 に切り出して
// usecase 層 や domain 層 の 業務ロジック には 関与させないようにする
type RefreshTokenIssuer interface {
	// リフレッシュトークン を 発行する
	// error を UsecaseError にしなかった理由は 依存先の実装(戻り値)に usecase 層の型を強制しないため
	// 実装が返した error は usecase 層で 文脈をつけて UsecaseError に変換する
	// よって interface としては 汎用的にしておく
	Issue() (*RefreshToken, error)
}

// リフレッシュトークンのハッシュ化
// interface にしている理由は リフレッシュトークン の ハッシュ化方法 を 柔軟にしておくため
// ハッシュ化方法 は 技術的な概念 になるので infra 層 に切り出して
// usecase 層 や domain 層 の 業務ロジック には 関与させないようにする
type RefreshTokenHasher interface {
	// ハッシュ化された リフレッシュトークン を 返す
	Hash(token string) (string, error)
}

// リフレッシュトークンの発行
// error を UsecaseError にした理由は このメソッドは 他の層で使われて
// 外部へ返すエラー契約を UsecaseError に 寄せるため
// 失敗原因が domain 層でも usecase 層 で包んで返す
func NewRefreshToken(token string, expiresAt time.Time) (*RefreshToken, UsecaseError) {
	d, err := domain.ReNewDate(expiresAt)
	if err != nil {
		return nil, WrapUsecaseError("リフレッシュトークンの作成に失敗しました", err)
	}
	return &RefreshToken{
		token:     token,
		expiresAt: d,
	}, nil
}

// 以下は ゲッター

func (r *RefreshToken) Token() string {
	return r.token
}

func (r *RefreshToken) ExpiresAt() domain.Date {
	return r.expiresAt
}
