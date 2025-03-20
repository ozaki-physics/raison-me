package usecase

import "github.com/ozaki-physics/raison-me/info/authN/domain"

// プレゼン層 との値のやり取り
type UserDto struct {
	accountID string
	userID    string
	userName  string
	pass      string
}

func constructorUserDto(aID string, uID string, uName string, pass string) *UserDto {
	// バリデーションは必要なのか?
	// ただ詰め替えるだけで 必要な値が入っているかのチェックは値オブジェクトの生成に任せる
	ud := &UserDto{
		aID,
		uID,
		uName,
		pass,
	}
	return ud
}

func NewUserDto(aID domain.AccountID, uID domain.UserID, uName domain.UserName, password domain.Password) *UserDto {
	return constructorUserDto(aID.Val(), uID.Val(), uName.Val(), password.String())
}

// プリミティブ型から生成するときは ReNew にする
func ReNewUserDto(aID string, uID string, uName string, pass string) *UserDto {
	return constructorUserDto(aID, uID, uName, pass)
}

// 以下 ゲッター

func (u *UserDto) AccountID() string {
	return u.accountID
}

func (u *UserDto) ID() string {
	return u.userID
}

func (u *UserDto) Name() string {
	return u.userName
}

// プレゼン層 から ユースケース層に渡すときは ユースケース層　で パスワード を取り出す必要があるが
// プレゼン層 で パスワードを取り出すことは無いと思うので 意図的にパブリックなゲッターを用意しない
func (u *UserDto) password() string {
	return u.pass
}
