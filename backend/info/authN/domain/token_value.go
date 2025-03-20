package domain

// トークン 値オブジェクト
type Token string

func NewToken(token string) (Token, error) {
	// TODO: 値としてのバリデーションが未実装
	t := Token(token)
	return t, nil
}

func ReNewToken(token string) (Token, DomainError) {
	t := Token(token)
	return t, nil
}
