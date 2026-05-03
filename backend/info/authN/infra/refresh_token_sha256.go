package infra

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/ozaki-physics/raison-me/info/authN/usecase"
)

type refreshTokenSHA256 struct {
	pepper string
}

func NewRefreshTokenSHA256(pepper string) (usecase.RefreshTokenHasher, error) {
	pepper = strings.TrimSpace(pepper)
	if pepper == "" {
		return nil, NewInfraError("リフレッシュトークンのペッパー は 必須です")
	}

	return &refreshTokenSHA256{pepper: pepper}, nil
}

func (h *refreshTokenSHA256) Hash(token string) (string, error) {
	if strings.TrimSpace(token) == "" {
		return "", NewInfraError("リフレッシュトークン が 空です")
	}

	sum := sha256.Sum256([]byte(h.pepper + token))
	return hex.EncodeToString(sum[:]), nil
}
