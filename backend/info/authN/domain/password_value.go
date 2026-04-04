package domain

import "golang.org/x/crypto/bcrypt"

// パスワード 値オブジェクト
type Password struct {
	hashedText string // ハッシュ化したパスワード
}

// 保存済みハッシュ文字列からパスワードオブジェクトを生成する
func NewPassword(passHash string) (Password, DomainError) {
	if passHash == "" {
		return Password{""}, NewDomainError("パスワードを入力してください")
	}

	if _, err := bcrypt.Cost([]byte(passHash)); err != nil {
		return Password{""}, WrapDomainError("パスワードハッシュが不正です", err)
	}

	return Password{passHash}, nil
}

// ハッシュ化されたパスワード から パスワードオブジェクトを生成する
func ReNewPassword(passHash string) (Password, DomainError) {
	return NewPassword(passHash)
}

// 以下ゲッター

func (p Password) HashedText() string {
	return p.hashedText
}

// 機密フィールドの出力形式変更
func (p Password) String() string {
	return "xxxxxx"
}

// 機密フィールドの出力形式変更
func (p Password) GoString() string {
	return "xxxxxx"
}
