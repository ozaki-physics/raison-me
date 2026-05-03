package presen

// アクセストークン だけ 返す レスポンスの JSON の 形式
type accessTokenResponse struct {
	AccessToken string `json:"accessToken"`
}

func NewAccessTokenResponse(accessToken string) *accessTokenResponse {
	return &accessTokenResponse{
		AccessToken: accessToken,
	}
}
