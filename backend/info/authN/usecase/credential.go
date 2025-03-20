package usecase

import "github.com/ozaki-physics/raison-me/info/authN/domain"

type Credential interface {
	// IDトークン の 発行
	Generate(accountID domain.AccountID, issuer string) (*domain.Credential, error)
	// IDトークン が 有効か
	IsAuthN(token domain.Token) (bool, error)
	// IDトークン の 無効化
	// Disable(id int) (bool, error)
}

type credential struct {
	credentialRepo domain.CredentialRepo
	cryptoRepo     domain.CryptoRepo
}

func NewCredential(credentialRepo domain.CredentialRepo, cryptoRepo domain.CryptoRepo) (Credential, error) {
	c := &credential{credentialRepo, cryptoRepo}
	return c, nil
}

func (c *credential) Generate(accountID domain.AccountID, issuer string) (*domain.Credential, error) {
	// TODO: 未実装
	// ストレージに保存
	return nil, nil
}

func (c *credential) IsAuthN(token domain.Token) (bool, error) {
	// TODO: 未実装
	// token から ユーザーを特定するもの と パスワードの代わりのトークンを取得
	// ユーザーを特定するもの で 紐づいてる クレデンシャルを取得
	// その中に一致するものがあるか確認
	return true, nil
}
