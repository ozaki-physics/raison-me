package cryptoasset

// CredentialCmc CMC と繋ぐためのクレデンシャルのインタフェース
type CredentialCmc interface {
	BaseUrl() string
	ApiKey() string
}

const (
	sandboxAPIKey       = "b54bcf4d-1bca-4e8e-9a24-22ff2c3d462c"
	sandboxBaseEndpoint = "https://sandbox-api.coinmarketcap.com"
	liveBaseEndpoint    = "https://pro-api.coinmarketcap.com"
)

// cmcDto CoinMarketCap の Credential の実体
type cmcDto struct {
	baseUrl string
	apiKey  string
}

func (cd *cmcDto) BaseUrl() string {
	return cd.baseUrl
}

func (cd *cmcDto) ApiKey() string {
	return cd.apiKey
}
