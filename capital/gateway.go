package capital

import (
	"log"
	"net/http"

	"github.com/ozaki-physics/raison-me/capital/infrastructure"
	"github.com/ozaki-physics/raison-me/capital/presentation"
	"github.com/ozaki-physics/raison-me/capital/usecase"
)

const serviceUrl = "/capital"

func MainCapital() {
	coinRepo := infrastructure.CreateCoinRepository()
	transactionRepo := infrastructure.CreateTransactionRepository()
	cryptoAssetsUsecase := usecase.CreateCryptoAssetsUsecase(coinRepo, transactionRepo)
	cryptoAssetsPresen := presentation.CreateCrypocryptoAssetsPresen(cryptoAssetsUsecase)
	// プレゼンテーション層の出力が プレゼンテーション層の入力になって 同じ層内で依存してて気持ちよくない
	cryptoAssetsRoutes := presentation.CreateCryptoAssetsRoutes(cryptoAssetsPresen)

	// ルーティングの定義
	http.HandleFunc(serviceUrl+"/crypto-assets/price", cryptoAssetsRoutes.Handler)
	// サーバ起動
	log.Println("Server Start >> http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
