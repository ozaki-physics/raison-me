package domain

import "time"

// リフレッシュトークン のクレームを表す
type RefreshTokenClaims struct {
	refreshTokenID RefreshTokenID
	accountID      AccountID
	tokenHash      string
	expiresAt      Date
	revokedAt      *Date
	createdAt      Date
	updatedAt      Date
}

func constructorRefreshTokenClaims(
	refreshTokenID RefreshTokenID,
	accountID AccountID,
	tokenHash string,
	expiresAt Date,
	revokedAt *Date,
	createdAt Date,
	updatedAt Date,
) (*RefreshTokenClaims, DomainError) {
	// TODO: バリデーションが十分か検討
	if refreshTokenID.IsNilID() {
		return nil, NewDomainError("refresh token ID が 不正です")
	}
	if accountID.IsNilID() {
		return nil, NewDomainError("account ID が 不正です")
	}
	if tokenHash == "" {
		return nil, NewDomainError("token hash が 不正です")
	}
	if expiresAt.IsZero() {
		return nil, NewDomainError("expiresAt が 不正です")
	}
	if createdAt.IsZero() {
		return nil, NewDomainError("createdAt が 不正です")
	}
	if updatedAt.IsZero() {
		return nil, NewDomainError("updatedAt が 不正です")
	}

	t := &RefreshTokenClaims{
		refreshTokenID: refreshTokenID,
		accountID:      accountID,
		tokenHash:      tokenHash,
		expiresAt:      expiresAt,
		revokedAt:      revokedAt,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
	return t, nil
}

// リフレッシュトークンの作成
// revokedAt は 新しい リフレッシュトークン を 作るときは nil になるから 引数に取らない
func NewRefreshTokenClaims(
	refreshTokenID RefreshTokenID,
	accountID AccountID,
	tokenHash string,
	expiresAt Date,
	now Date,
) (*RefreshTokenClaims, DomainError) {
	return constructorRefreshTokenClaims(refreshTokenID, accountID, tokenHash, expiresAt, nil, now, now)
}

func ReNewRefreshTokenClaims(
	refreshTokenID string,
	accountID string,
	tokenHash string,
	expiresAt time.Time,
	revokedAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) (*RefreshTokenClaims, DomainError) {
	var errs []DomainError

	rID, err01 := ReNewRefreshTokenID(refreshTokenID)
	errs = append(errs, err01)
	aID, err02 := ReNewAccountID(accountID)
	errs = append(errs, err02)

	// tokenHash は ドメインオブジェクトの生成に必要な値だけど 値オブジェクト にはしない
	th := tokenHash

	at, err03 := ReNewDate(expiresAt)
	errs = append(errs, err03)
	// revokedAt は 意図的に 失効 していないと nil になる
	// よって 値オブジェクト を生成する前に nil チェックする
	var revokedAtPtr *Date
	if revokedAt != nil {
		revokedAtVal, err04 := ReNewDate(*revokedAt)
		errs = append(errs, err04)
		if err04 == nil {
			revokedAtPtr = &revokedAtVal
		}
	}
	cAt, err05 := ReNewDate(createdAt)
	errs = append(errs, err05)
	uAt, err06 := ReNewDate(updatedAt)
	errs = append(errs, err06)

	// TODO: 最初に発生したエラー以外も検知できるようにしたい
	for _, v := range errs {
		if v != nil {
			return nil, v
		}
	}

	return constructorRefreshTokenClaims(rID, aID, th, at, revokedAtPtr, cAt, uAt)
}

// リフレッシュトークンが 失効しているか
func (r *RefreshTokenClaims) IsRevoked() bool {
	return r.revokedAt != nil
}

// リフレッシュトークンが 有効期限切れか
func (r *RefreshTokenClaims) IsExpired(now Date) bool {
	return !r.expiresAt.After(now.Time)
}

// リフレッシュトークンが 有効か
func (r *RefreshTokenClaims) IsActive(now Date) bool {
	return !r.IsRevoked() && !r.IsExpired(now)
}

// ローテーションした 新しい リフレッシュトークン オブジェクト を生成する
func (r *RefreshTokenClaims) WithRotatedToken(tokenHash string, expiresAt Date, now Date) (*RefreshTokenClaims, DomainError) {
	if tokenHash == "" {
		return nil, NewDomainError("token hash が 不正です")
	}
	if expiresAt.IsZero() {
		return nil, NewDomainError("expiresAt が 不正です")
	}

	return constructorRefreshTokenClaims(r.refreshTokenID, r.accountID, tokenHash, expiresAt, r.revokedAt, r.createdAt, now)
}

// リフレッシュトークンを 失効させた リフレッシュトークン オブジェクト を生成する
func (r *RefreshTokenClaims) WithRevoked(now Date) (*RefreshTokenClaims, DomainError) {
	if now.IsZero() {
		return nil, NewDomainError("revokedAt が 不正です")
	}

	return constructorRefreshTokenClaims(r.refreshTokenID, r.accountID, r.tokenHash, r.expiresAt, &now, r.createdAt, now)
}

// 以下ゲッター

func (r *RefreshTokenClaims) ID() RefreshTokenID {
	return r.refreshTokenID
}

func (r *RefreshTokenClaims) AccountID() AccountID {
	return r.accountID
}

func (r *RefreshTokenClaims) TokenHash() string {
	return r.tokenHash
}

func (r *RefreshTokenClaims) ExpiresAt() Date {
	return r.expiresAt
}

func (r *RefreshTokenClaims) RevokedAt() *Date {
	return r.revokedAt
}

func (r *RefreshTokenClaims) CreatedAt() Date {
	return r.createdAt
}

func (r *RefreshTokenClaims) UpdatedAt() Date {
	return r.updatedAt
}
