package domain

import "golang.org/x/crypto/bcrypt"

// パスワード 値オブジェクト
type Password struct {
	pass string // ハッシュ化したパスワード
}

// 平文パスワード から ハッシュ化されたパスワードオブジェクトを生成する
func NewPassword(pass string) (Password, DomainError) {
	if pass == "" {
		return Password{""}, NewDomainError("パスワードを入力してください")
	}
	// TODO: 使える文字種の制限

	// ハッシュ化
	bPass := []byte(pass)
	// bcrypt の仕様で 72バイト以上が受け取れないから
	if len(bPass) >= 72 {
		return Password{""}, NewDomainError("パスワードを短くしてください")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword(bPass, bcrypt.DefaultCost)
	if err != nil {
		return Password{""}, WrapDomainError("パスワードがハッシュ化できませんでした", err)
	}

	return Password{string(hashedPassword)}, nil
}

// ハッシュ化されたパスワード から パスワードオブジェクトを生成する
func ReNewPassword(pass string) (Password, DomainError) {
	return Password{pass}, nil
}

// 保有してるハッシュ値と入力された値が一致しているか
func (p *Password) isSame(inputPassword string) (bool, DomainError) {
	err := bcrypt.CompareHashAndPassword([]byte(p.ToHash()), []byte(inputPassword))
	if err != nil {
		return false, NewDomainError("パスワードが一致しません")
	} else {
		return true, nil
	}
}

// 以下ゲッター

func (p *Password) ToHash() string {
	return p.pass
}

// 機密フィールドの出力形式変更
func (p Password) String() string {
	return "xxxxxx"
}

// 機密フィールドの出力形式変更
func (p Password) GoString() string {
	return "xxxxxx"
}
