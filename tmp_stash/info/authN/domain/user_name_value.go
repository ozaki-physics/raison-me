package domain

// ユーザー名 値オブジェクト
type UserName struct {
	name string
}

func NewUserName(name string) (UserName, DomainError) {
	if name == "" {
		return UserName{""}, NewDomainError("ユーザー名を入力してください")
	}
	// TODO: 使える文字種の制限
	// TODO: 文字数の制限

	u := UserName{name}
	return u, nil
}

// 以下ゲッター

func (u *UserName) Val() string {
	return u.name
}
