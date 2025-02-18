package domain

// あえてユーザーという名称を使う
type User struct {
	accountID AccountID // システム上の ID(ユーザー+管理者 で一意)
	id        UserID    // 機能上の ID(ユーザー界で一意)
	name      UserName
}

func constructorUser(aID AccountID, uID UserID, uName UserName) (*User, DomainError) {
	// TODO: 引数の値が ゼロ値じゃないかチェック
	// そもそもドメインオブジェクトが完全コンストラクタだから必要ない?
	if aID.IsNilID() {
		return nil, NewDomainError("引数のAccountIDが不正です")
	}

	u := &User{
		aID,
		uID,
		uName,
	}
	return u, nil
}

func NewUser(aID AccountID, uID UserID, uName UserName) (*User, DomainError) {
	return constructorUser(aID, uID, uName)
}

// インフラ層でドメインオブジェクトを生成するため
// 本当は 値オブジェクト を1個ずつ生成した方がいいんだろうけど 面倒だからまとめた
func ReNewUser(accountID string, userID string, name string) (*User, DomainError) {
	var errs []DomainError
	aID, err01 := ReNewAccountID(accountID)
	errs = append(errs, err01)
	uID, err02 := NewUserID(userID)
	errs = append(errs, err02)
	uName, err03 := NewUserName(name)
	errs = append(errs, err03)

	for _, v := range errs {
		if v != nil {
			return nil, v
		}
	}
	return constructorUser(aID, uID, uName)
}

// 以下ゲッター

func (u *User) AccountID() AccountID {
	return u.accountID
}

func (u *User) ID() UserID {
	return u.id
}

func (u *User) Name() UserName {
	return u.name
}
