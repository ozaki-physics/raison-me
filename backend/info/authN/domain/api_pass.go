package domain

type APIPass struct {
	APIKeyID  int
	accountID AccountID // User との紐づけ
	key       APIKey
	iat       Date // 発行日付
}
