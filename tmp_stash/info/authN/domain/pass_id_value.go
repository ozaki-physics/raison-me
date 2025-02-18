package domain

// パスワードID 値オブジェクト
type PassID struct {
	id
}

func constructorPassID(id id) (PassID, DomainError) {
	if id.prefixID != "p" {
		return PassID{NilID()}, NewDomainError("PassIDに設定できないプレフィックスです")
	}
	// TODO: ユーザー+管理者 で一意になるように 自動で採番

	return PassID{id}, nil
}

func NewPassID() (PassID, DomainError) {
	id, err := NewID("p")
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
