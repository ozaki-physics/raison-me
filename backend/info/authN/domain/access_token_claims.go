package domain

import "time"

// アクセストークン のクレームを表す
type AccessTokenClaims struct {
	AccountID AccountID
	IssuedAt  Date
	ExpiresAt Date
}

func constructorAccessTokenClaims(
	accountID AccountID,
	issuedAt Date,
	expiresAt Date,
) (*AccessTokenClaims, DomainError) {
	// TODO: バリデーションが十分か検討
	if accountID.IsNilID() {
		return nil, NewDomainError("account ID が 不正です")
	}
	if issuedAt.IsZero() {
		return nil, NewDomainError("issuedAt が 不正です")
	}
	if expiresAt.IsZero() {
		return nil, NewDomainError("expiresAt が 不正です")
	}

	t := &AccessTokenClaims{
		AccountID: accountID,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}
	return t, nil
}

func NewAccessTokenClaims(
	accountID AccountID,
	issuedAt Date,
	expiresAt Date,
) (*AccessTokenClaims, DomainError) {
	return constructorAccessTokenClaims(accountID, issuedAt, expiresAt)
}

func ReNewAccessTokenClaims(
	accountID string,
	issuedAt time.Time,
	expiresAt time.Time,
) (*AccessTokenClaims, DomainError) {
	var errs []DomainError

	aID, err := ReNewAccountID(accountID)
	errs = append(errs, err)
	iAt, err := ReNewDate(issuedAt)
	errs = append(errs, err)
	eAt, err := ReNewDate(expiresAt)
	errs = append(errs, err)

	// TODO: 最初に発生したエラー以外も検知できるようにしたい
	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}
	a, err := constructorAccessTokenClaims(aID, iAt, eAt)
	return a, err
}
