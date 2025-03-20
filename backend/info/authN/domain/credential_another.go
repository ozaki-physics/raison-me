package domain

type Credential02 struct {
	id             int
	accountID      AccountID // ユーザーと紐づけるため
	token          Token     // トークン
	expirationTime Date      // 有効期限
	issuer         string    // 発行者
	issuedAt       Date      // 発行日時
}

func NewCredential02(id int, accountID AccountID, token Token, expirationTime Date, issuer string, issuedAt Date) (*Credential02, DomainError) {
	c := &Credential02{
		id:             id,
		accountID:      accountID,
		token:          token,
		expirationTime: expirationTime,
		issuer:         issuer,
		issuedAt:       issuedAt,
	}
	return c, nil
}

func PrimitiveNewCredential(id int, accountID string, token string, expirationTime string, issuer string, issuedAt string) (*Credential02, DomainError) {
	// TODO: ドメインとしてのバリデーション が未実装

	aid, err := ReNewAccountID(accountID)
	if err != nil {
		return nil, err
	}
	t, err := ReNewToken(token)
	if err != nil {
		return nil, err
	}
	exp, err := ReNewDate(expirationTime)
	if err != nil {
		return nil, err
	}
	iat, err := ReNewDate(issuedAt)
	if err != nil {
		return nil, err
	}

	c, err := NewCredential02(id, aid, t, exp, issuer, iat)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Credential02) ID() int {
	return c.id
}

func (c *Credential02) AccountID() AccountID {
	return c.accountID
}

func (c *Credential02) Token() Token {
	return c.token
}

func (c *Credential02) ExpirationTime() Date {
	return c.expirationTime
}
func (c *Credential02) Issuer() string {
	return c.issuer
}
func (c *Credential02) IssuedAt() Date {
	return c.issuedAt
}

func (c *Credential02) IsAuthN(token Token) (bool, error) {
	if c.Token() == token {
		return true, nil
	}
	return false, nil
}
