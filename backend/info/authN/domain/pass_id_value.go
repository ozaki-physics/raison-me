package domain

// パスワードID 値オブジェクト
type PassID struct {
	id
}

func constructorPassID(id id) (PassID, DomainError) {
	if id.IsNilID() {
		return PassID{NilID()}, NewDomainError("PassIDが不正です")
	}

	// ユーザー+管理者 で一意にしたい
	// UUID v7 なら衝突の可能性は 限りなく低いので チェックはしない

	return PassID{id}, nil
}

func NewPassID() (PassID, DomainError) {
	id, err := NewID()
	if err != nil {
		return PassID{NilID()}, err
	}

	return constructorPassID(id)
}

func ReNewPassID(data string) (PassID, DomainError) {
	id, err := ReNewID(data)
	if err != nil {
		return PassID{NilID()}, err
	}
	return constructorPassID(id)
}

// 以下ゲッター
