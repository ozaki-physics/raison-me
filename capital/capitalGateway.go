// Service : capital のコンテキストたちの DI などを行う
package capital

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	infra "github.com/ozaki-physics/raison-me/capital/infrastructure/cryptoAsset"
	share "github.com/ozaki-physics/raison-me/capital/infrastructure/share"
	presen "github.com/ozaki-physics/raison-me/capital/presentation/cryptoAsset"
	usecase "github.com/ozaki-physics/raison-me/capital/usecase/cryptoAsset"
	"github.com/ozaki-physics/raison-me/share/config"
)

func CryptoAsset() chi.Router {
	cmcCredential := infra.CreateCredentialCoinMarketCapGcp(config.IsLive)
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
	r := chi.NewRouter()
	r.Route("/crypto-assets", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("これは capital の cryptoAssets だよ\n"))
		})

		r.HandleFunc("/price", apiHandler.Handler)
		if config.IsLive {
			r.HandleFunc("/line", lineController.SoundReflection)
		}
	})
	return r
}
