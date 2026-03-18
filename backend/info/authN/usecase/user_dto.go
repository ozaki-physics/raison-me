package usecase

import "github.com/ozaki-physics/raison-me/info/authN/domain"

// ユースケース層 から プレゼン層 へ 連携
// 値を詰め替えるための構造体 だから DTO(Data Transfer Object) と呼ぶ
// 理由は 業務ロジック が プレゼン層 に 漏れないようにするため
// 業務ロジック を 持たないため 値は 公開して良い
// パスワードは プレゼン層 で 取り出すことは無いと思うので 存在させない
type UserDto struct {
	AccountID string
	UserID    string
	UserName  string
}

func constructorUserDto(aID string, uID string, uName string) *UserDto {
	// バリデーションは必要なのか?
	// ただ詰め替えるだけで 必要な値が入っているかのチェックは値オブジェクトの生成に任せる
	ud := &UserDto{
		aID,
		uID,
		uName,
	}
	return ud
}

func NewUserDto(aID domain.AccountID, uID domain.UserID, uName domain.UserName) *UserDto {
	return constructorUserDto(aID.Val(), uID.Val(), uName.Val())
}

// プリミティブ型から生成するときは ReNew にする
func ReNewUserDto(aID string, uID string, uName string) *UserDto {
	return constructorUserDto(aID, uID, uName)
}
