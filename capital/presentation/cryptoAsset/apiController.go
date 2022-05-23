package cryptoasset

import (
	"encoding/json"
	"fmt"
	"net/http"

	usecase "github.com/ozaki-physics/raison-me/capital/usecase/cryptoAsset"
)

type ApiController interface {
	Get(w http.ResponseWriter, r *http.Request) error
	Post(w http.ResponseWriter, r *http.Request) error
	// Put(w http.ResponseWriter, r *http.Request) error
	// Delete(w http.ResponseWriter, r *http.Request) error
}

type apiController struct {
	caUsecase usecase.CryptoAssetsUsecase
}

func CreateApiController(u usecase.CryptoAssetsUsecase) ApiController {
	return &apiController{u}
}

func (apic *apiController) Get(w http.ResponseWriter, r *http.Request) error {
	// URL パラメータ を string で取得
	paramSymbol := r.URL.Query().Get("symbol")

	averagePriceRate, err := apic.caUsecase.CoinAveragePrice(paramSymbol)
	if err != nil {
		return err
	}

	type response struct {
		Symbol           string  `json:"symbol"`
		AveragePriceRate float64 `json:"price"`
	}

	res := response{
		paramSymbol,
		averagePriceRate,
	}
	// http レスポンスに格納する
	if err = json.NewEncoder(w).Encode(res); err != nil {
		return err
	}
	return nil
}

func (apic *apiController) Post(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "Hello, World!\n")
	return nil
}
