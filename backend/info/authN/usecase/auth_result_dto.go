package usecase

// 認証の結果 API で結果を保持する DTO
type AuthResultDto struct {
	AccessToken  string
	RefreshToken string
}

func NewAuthResultDto(accessToken string, refreshToken string) *AuthResultDto {
	return &AuthResultDto{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
