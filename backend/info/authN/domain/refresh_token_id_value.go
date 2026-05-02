package domain

// RefreshTokenID 値オブジェクト
type RefreshTokenID struct {
	id
}

func constructorRefreshTokenID(id id) (RefreshTokenID, DomainError) {
	if id.IsNilID() {
		return RefreshTokenID{NilID()}, NewDomainError("refresh token ID が 不正です")
	}

	return RefreshTokenID{id}, nil
}

func NewRefreshTokenID() (RefreshTokenID, DomainError) {
	id, err := NewID()
	if err != nil {
		return RefreshTokenID{NilID()}, err
	}

	return constructorRefreshTokenID(id)
}

func ReNewRefreshTokenID(data string) (RefreshTokenID, DomainError) {
	id, err := ReNewID(data)
	if err != nil {
		return RefreshTokenID{NilID()}, err
	}

	return constructorRefreshTokenID(id)
}
