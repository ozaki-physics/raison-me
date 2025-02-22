package cryptoasset

import share "github.com/ozaki-physics/raison-me/capital/infrastructure/share"

// dataCoinMarketCap JSON から struct に変換する
type dataCoinMarketCap struct {
	Data struct {
		Service struct {
			ApiKey string `json:"apiKey"`
		} `json:"CoinMarketCap"`
	} `json:"data"`
}

func CreateCredentialCoinMarketCapJson(isLive bool) CredentialCmc {
	// Sandbox モード
	if isLive == false {
		return &cmcDto{sandboxBaseEndpoint, sandboxAPIKey}
	}

	// Live モード
	var d dataCoinMarketCap
	share.JsonToStruct(d, "./capital/infrastructure/cryptoAsset/json/key.json")
	return &cmcDto{liveBaseEndpoint, d.Data.Service.ApiKey}

}
