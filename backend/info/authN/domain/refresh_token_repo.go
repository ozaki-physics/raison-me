package domain

import "context"

// この struct は ユースケース層 で使われていて 実装は インフラ層
type RefreshTokenRepo interface {
	// リフレッシュトークン を 登録
	Insert(ctx context.Context, refreshToken *RefreshTokenClaims) (*RefreshTokenClaims, error)
	// トークンハッシュ から リフレッシュトークン を 取得
	FindByTokenHash(ctx context.Context, tokenHash string) (*RefreshTokenClaims, error)
	// リフレッシュトークン を ローテーション
	Rotate(ctx context.Context, currentTokenHash string, nextRefreshToken *RefreshTokenClaims) (*RefreshTokenClaims, error)
	// リフレッシュトークン を 失効
	Revoke(ctx context.Context, tokenHash string) error
}
