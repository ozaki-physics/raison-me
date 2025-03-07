package cryptoasset

import (
	"log"

	share "github.com/ozaki-physics/raison-me/capital/infrastructure/share"
)

func CreateCredentialCoinMarketCapGcp(isLive bool) CredentialCmc {
	// Sandbox モード
	if !isLive {
		return &cmcDto{sandboxBaseEndpoint, sandboxAPIKey}
	}

	// Live モード
	var cd cmcDto
	cd.baseUrl = liveBaseEndpoint
	value, err := share.GcpSecretValue("CoinMarketCap_API", 1)
	if err != nil {
		log.Fatalln(err)
	}
	cd.apiKey = value

	return &cd
}
