package infra

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/ozaki-physics/raison-me/info/authN/usecase"
)

// リフレッシュトークン の 有効期限 (7日)
const refreshTokenTTLSeconds = 604800

type refreshTokenRandom struct{}

func NewRefreshTokenRandom() usecase.RefreshTokenIssuer {
	return &refreshTokenRandom{}
}

// ランダムなバイト列を生成してリフレッシュトークンとする実装
func (g *refreshTokenRandom) Issue() (*usecase.RefreshToken, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return nil, WrapInfraError("リフレッシュトークン生成エラー: %w", err)
	}
	t := base64.RawURLEncoding.EncodeToString(buf)

	now := time.Now().UTC()
	exp := now.Add(time.Duration(refreshTokenTTLSeconds) * time.Second)

	rtd, err := usecase.NewRefreshToken(t, exp)
	if err != nil {
		return nil, WrapInfraError("リフレッシュトークン生成エラー: %w", err)
	}
	return rtd, nil
}
