package domain

// アカウントID 値オブジェクト
type AccountID struct {
	id
}

func constructorAccountID(id id) (AccountID, DomainError) {
	if id.IsNilID() {
		return AccountID{NilID()}, NewDomainError("AccountIDが不正です")
	}

	// ユーザー+管理者 で一意にしたい
	// UUID v7 なら衝突の可能性は 限りなく低いので チェックはしない

	return AccountID{id}, nil
}

func NewAccountID() (AccountID, DomainError) {
	id, err := NewID()
	if err != nil {
		return AccountID{NilID()}, err
	}

	return constructorAccountID(id)
}

func ReNewAccountID(data string) (AccountID, DomainError) {
	id, err := ReNewID(data)
	if err != nil {
		return AccountID{NilID()}, err
	}
	return constructorAccountID(id)
}

// 以下ゲッター
