package domain

// パスワード
type Pass struct {
	passID    PassID    // システム上の ID
	accountID AccountID // User との紐づけ
	password  Password  // ハッシュ化されたパスワード
	iat       Date      // 発行日付
}

func constructorPass(pID PassID, aID AccountID, pass Password, iat Date) (*Pass, DomainError) {
	// TODO: 引数の値が ゼロ値じゃないかチェック
	// そもそもドメインオブジェクトが完全コンストラクタだから必要ない?
	if pID.IsNilID() {
		return nil, NewDomainError("引数のPassIDが不正です")
	}
	if aID.IsNilID() {
		return nil, NewDomainError("引数のAccountIDが不正です")
	}

	p := &Pass{
		passID:    pID,
		accountID: aID,
		password:  pass,
		iat:       iat,
	}
	return p, nil
}

func NewPass(pID PassID, aID AccountID, pass Password, iat Date) (*Pass, DomainError) {
	return constructorPass(pID, aID, pass, iat)
}

// インフラ層でドメインオブジェクトを生成するため
// 本当は 値オブジェクト を1個ずつ生成した方がいいんだろうけど 面倒だからまとめた
func ReNewPass(passID string, accountID string, password string, iat string) (*Pass, DomainError) {
	var errs []DomainError
	pID, err01 := ReNewPassID(passID)
	errs = append(errs, err01)
	aID, err02 := ReNewAccountID(accountID)
	errs = append(errs, err02)
	p, err03 := ReNewPassword(password)
	errs = append(errs, err03)
	at, err04 := ReNewDate(iat)
	errs = append(errs, err04)

	// TODO: 最初に発生したエラー以外も検知できるようにしたい
	for _, v := range errs {
		if v != nil {
			return nil, v
		}
	}
	return constructorPass(pID, aID, p, at)
}

func (p *Pass) IsLogin(inputPassword string) (bool, DomainError) {
	return p.password.isSame(inputPassword)
}

// 以下ゲッター

func (p *Pass) PassID() PassID {
	return p.passID
}

func (p *Pass) AccountID() AccountID {
	return p.accountID
}

func (p *Pass) Password() Password {
	return p.password
}

func (p *Pass) IssuedAt() Date {
	return p.iat
}
