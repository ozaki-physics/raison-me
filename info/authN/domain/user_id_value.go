package domain

// ユーザーID 値オブジェクト
type UserID struct {
	id string
}

func NewUserID(id string) (UserID, DomainError) {
	if id == "" {
		return UserID{""}, NewDomainError("ユーザーIDを入力してください")
	}
	// TODO: 使える文字種の制限
	// TODO: ユーザー界で一意のチェック
	// TODO: 文字数の制限

	u := UserID{id}
	return u, nil
}

// 以下ゲッター

func (u *UserID) Val() string {
	return u.id
}
