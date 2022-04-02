package infrastructure

import (
	"log"

	"github.com/ozaki-physics/raison-me/config"
)

const (
	sandboxAPIKey       = "b54bcf4d-1bca-4e8e-9a24-22ff2c3d462c"
	sandboxBaseEndpoint = "https://sandbox-api.coinmarketcap.com"
	liveBaseEndpoint    = "https://pro-api.coinmarketcap.com"
)

// coinMarketCap クレデンシャル(key と URL を保持)
type coinMarketCap struct {
	Key     string
	baseURL string
}

// responseStatus レスポンスの共通部分
type responseStatus struct {
	Timestamp    string `json:"timestamp"`
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Elapsed      int    `json:"elapsed"`
	CreditCount  int    `json:"credit_count"`
	Notice       string `json:"notice"`
}

// platform CMCIDMap と Metadata の共通部分
type platform struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	Slug         string `json:"slug"`
	TokenAddress string `json:"token_address"`
}

// getCredential CoinMarketCap にアクセスするクレデンシャルを保持
func getCredential() coinMarketCap {
	// Sandbox モード
	if config.IsLive == false {
		return coinMarketCap{
			sandboxAPIKey,
			sandboxBaseEndpoint,
		}
	}

	// Live モード
	CMCAPIKey, err := config.GetGCPSecretValue("CoinMarketCap_API", 1)
	if err != nil {
		log.Fatalln(err)
	}
	return coinMarketCap{
		CMCAPIKey,
		liveBaseEndpoint,
	}
}
