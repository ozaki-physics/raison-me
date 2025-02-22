package domain

// アカウントID 値オブジェクト
type AccountID struct {
	id
}

func constructorAccountID(id id) (AccountID, DomainError) {
	if id.prefixID != "a" {
		return AccountID{NilID()}, NewDomainError("AccountIDに設定できないプレフィックスです")
	}
	// TODO: ユーザー+管理者 で一意になるように 自動で採番

	return AccountID{id}, nil
}

func NewAccountID() (AccountID, DomainError) {
	id, err := NewID("a")
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
