package infra

import (
	"strings"

	"github.com/ozaki-physics/raison-me/info/authN/domain"
	"golang.org/x/crypto/bcrypt"
)

type passwordBcrypt struct {
	pepper string
}

// 戻り値 が インタフェース だから 実装を強制できる
func NewPasswordBcrypt(pepper string) (domain.PasswordHasherRepo, error) {
	if strings.TrimSpace(pepper) == "" {
		return nil, NewInfraError("AUTHN_PASSWORD_PEPPER が設定されていません")
	}

	return passwordBcrypt{pepper: pepper}, nil
}

func (h passwordBcrypt) Hash(plain string) (string, error) {
	if plain == "" {
		return "", NewInfraError("パスワードを入力してください")
	}

	passByte := []byte(h.pepper + plain)
	if len(passByte) >= 72 {
		return "", NewInfraError("パスワードを短くしてください")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(passByte, bcrypt.DefaultCost)
	if err != nil {
		return "", WrapInfraError("パスワードがハッシュ化できませんでした", err)
	}

	return string(hashedPassword), nil
}

func (h passwordBcrypt) Verify(password domain.Password, plain string) error {
	hashed := password.HashedText()
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(h.pepper+plain)); err != nil {
		return NewInfraError("パスワードが一致しません")
	}

	return nil
}
