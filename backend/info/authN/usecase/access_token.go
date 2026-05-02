package usecase

import (
	"time"

	"github.com/ozaki-physics/raison-me/info/authN/domain"
)

// アクセストークンの発行結果
type AccessToken struct {
	token     string
	expiresAt domain.Date
}

// アクセストークンの発行
// (New メソッドの代わり)
// interface にしている理由は アクセストークン の 発行方法 を 柔軟にしておくため
// 発行方法 は 技術的な概念 になるので infra 層 に切り出して
// usecase 層 や domain 層 の 業務ロジック には 関与させないようにする
type AccessTokenIssuer interface {
	// accountID を 引数 に して アクセストークン を 発行する
	Issue(accountID domain.AccountID) (*AccessToken, error)
}

// アクセストークンの検証
// interface にしている理由は アクセストークン の 検証方法 を 柔軟にしておくため
// 検証方法 は 技術的な概念 になるので infra 層 に切り出して
// usecase 層 や domain 層 の 業務ロジック には 関与させないようにする
type AccessTokenVerifier interface {
	// token を 引数 に して アクセストークン を 検証する
	// 検証に成功した場合は アクセストークンのクレーム を 返す
	Verify(token string) (*domain.AccessTokenClaims, error)
}

func ReNewAccessToken(token string, exp time.Time) (*AccessToken, UsecaseError) {
	d, err := domain.ReNewDate(exp)
	if err != nil {
		return nil, WrapUsecaseError("failed to create date", err)
	}
	return &AccessToken{
		token:     token,
		expiresAt: d,
	}, nil
}

// 以下は ゲッター

func (a *AccessToken) Token() string {
	return a.token
}

func (a *AccessToken) ExpiresAt() domain.Date {
	return a.expiresAt
}
