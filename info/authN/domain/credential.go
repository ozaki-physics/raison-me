package domain

type Credential struct {
	user User
	pass Pass
}

func NewCredential(u User, p Pass) (*Credential, DomainError) {
	if &u == nil {
		return nil, NewDomainError("ユーザーオブジェクトが存在しません")
	}
	if &p == nil {
		return nil, NewDomainError("パスワードオブジェクトが存在しません")
	}

	c := &Credential{u, p}
	return c, nil
}

// 以下ゲッター

func (c *Credential) User() User {
	return c.user
}

func (c *Credential) Pass() Pass {
	return c.pass
}
