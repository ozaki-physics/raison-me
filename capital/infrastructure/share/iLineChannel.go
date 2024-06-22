package share

// CredentialLine LINE bot と繋ぐためのクレデンシャルのインタフェース
type CredentialLine interface {
	Secret() string
	Token() string
}

// cmcDto LINE Channel の Credential の実体
type lineDto struct {
	secret string
	token  string
}

func (ld *lineDto) Secret() string {
	return ld.secret
}

func (ld *lineDto) Token() string {
	return ld.token
}
