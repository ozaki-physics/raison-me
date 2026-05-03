package infra

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ozaki-physics/raison-me/info/authN/domain"
)

type refreshTokenRepoSQL struct {
	pool *pgxpool.Pool
}

// 戻り値 が インタフェース だから 実装を強制できる
func NewRefreshTokenRepoSQL(pool *pgxpool.Pool) (domain.RefreshTokenRepo, error) {
	return &refreshTokenRepoSQL{pool: pool}, nil
}

func (r *refreshTokenRepoSQL) Insert(ctx context.Context, refreshToken *domain.RefreshTokenClaims) (*domain.RefreshTokenClaims, error) {
	sqlStatement := `
		INSERT INTO app.refresh_tokens (
			  refresh_token_id
			, account_id
			, token_hash
			, expires_at
			, revoked_at
			, created_at
			, updated_at
		) VALUES (
			 $1
		 , $2
		 , $3
		 , $4
		 , $5
		 , $6
		 , $7
		)
	`

	rID := refreshToken.ID()
	aID := refreshToken.AccountID()
	// revokedAt は NULL になる可能性があるため ポインタで扱う
	var revokedAt *time.Time
	if refreshToken.RevokedAt() != nil {
		t := refreshToken.RevokedAt().Time
		revokedAt = &t
	}

	_, err := r.pool.Exec(
		ctx,
		sqlStatement,
		rID.Val(),
		aID.Val(),
		refreshToken.TokenHash(),
		refreshToken.ExpiresAt().Time,
		revokedAt,
		refreshToken.CreatedAt().Time,
		refreshToken.UpdatedAt().Time,
	)
	if err != nil {
		log.Printf("exec: %v\n", err)
		return nil, err
	}

	return refreshToken, nil
}

func (r *refreshTokenRepoSQL) FindByTokenHash(ctx context.Context, tokenHash string) (*domain.RefreshTokenClaims, error) {
	sqlStatement := `
		SELECT
			  refresh_token_id
			, account_id
			, token_hash
			, expires_at
			, revoked_at
			, created_at
			, updated_at
		FROM app.refresh_tokens
		WHERE token_hash = $1
		LIMIT 1
	`

	row := r.pool.QueryRow(ctx, sqlStatement, tokenHash)
	return scanRefreshToken(row)
}

func (r *refreshTokenRepoSQL) Rotate(ctx context.Context, currentTokenHash string, nextRefreshToken *domain.RefreshTokenClaims) (*domain.RefreshTokenClaims, error) {
	sqlStatement := `
		UPDATE app.refresh_tokens
		SET
			  token_hash = $2
			, expires_at = $3
			, updated_at = $4
		WHERE
					token_hash = $1
			AND revoked_at IS NULL
			AND expires_at > $4
		RETURNING
			  refresh_token_id
			, account_id
			, token_hash
			, expires_at
			, revoked_at
			, created_at
			, updated_at
	`

	row := r.pool.QueryRow(
		ctx,
		sqlStatement,
		currentTokenHash,
		nextRefreshToken.TokenHash(),
		nextRefreshToken.ExpiresAt().Time,
		nextRefreshToken.UpdatedAt().Time,
	)
	return scanRefreshToken(row)
}

func (r *refreshTokenRepoSQL) Revoke(ctx context.Context, tokenHash string) error {
	sqlStatement := `
		UPDATE app.refresh_tokens
		SET
			  revoked_at = $2
			, updated_at = $3
		WHERE
					token_hash = $1
			AND revoked_at IS NULL
	`

	commandTag, err := r.pool.Exec(
		ctx,
		sqlStatement,
		tokenHash,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		log.Printf("exec: %v\n", err)
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// pgx.Row から domain.RefreshToken を構築するヘルパー関数
func scanRefreshToken(row pgx.Row) (*domain.RefreshTokenClaims, error) {
	// DB の 値を保持するための変数
	var refreshTokenID, accountID, tokenHash string
	var expiresAt, createdAt, updatedAt time.Time
	// DB の NULL 値を扱うために sql.NullTime を使用
	var revokedAt sql.NullTime

	// DB から値をスキャン
	err := row.Scan(&refreshTokenID, &accountID, &tokenHash, &expiresAt, &revokedAt, &createdAt, &updatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.NewDomainError("リフレッシュトークンが見つかりません")
		}
		log.Printf("scan: %v\n", err)
		return nil, err
	}

	// sql.NullTime を *time.Time に変換
	var revokedAtPtr *time.Time
	if revokedAt.Valid {
		revokedAtVal := revokedAt.Time
		revokedAtPtr = &revokedAtVal
	}

	refreshToken, err := domain.ReNewRefreshTokenClaims(refreshTokenID, accountID, tokenHash, expiresAt, revokedAtPtr, createdAt, updatedAt)
	if err != nil {
		log.Printf("renew refresh token: %v\n", err)
		return nil, err
	}

	return refreshToken, nil
}
