// Service : capital のコンテキストたちの DI などを行う
package capital

import (
	"log"
	"net/http"
	"os"

	infra "github.com/ozaki-physics/raison-me/capital/infrastructure/cryptoAsset"
	share "github.com/ozaki-physics/raison-me/capital/infrastructure/share"
	presen "github.com/ozaki-physics/raison-me/capital/presentation/cryptoAsset"
	usecase "github.com/ozaki-physics/raison-me/capital/usecase/cryptoAsset"
)

const serviceUrl = "/capital"

func CryptoAsset() {
	// cmcCredential := infra.CreateCredentialCoinMarketCapJson(false)
	cmcCredential := infra.CreateCredentialCoinMarketCapGcp(true)
	cmcIds := infra.CreateCMCIdsJson()
	coinRepo := infra.CreateCoinRepository(cmcCredential, cmcIds)
	transactionRepo := infra.CreateTransactionRepository()
	cryptoAssetUsecase := usecase.CreateCryptoAssetUsecase(coinRepo, transactionRepo)
	apiController := presen.CreateApiController(cryptoAssetUsecase)
	// REST API にするために
	apiHandler := presen.CreateApiHandler(apiController)

	lineCredential := share.CreateCredentialLineGcp()
	lineController := presen.CreateLineController(lineCredential, cryptoAssetUsecase)

	// ルーティングの定義
	http.HandleFunc(serviceUrl+"/crypto-assets/price", apiHandler.Handler)
	http.HandleFunc(serviceUrl+"/crypto-assets/line", lineController.SoundReflection)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)

	// サーバ起動
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
