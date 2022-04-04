package presentation

import (
	"encoding/json"
	"net/http"

	"github.com/ozaki-physics/raison-me/capital/usecase"
)

// CryptoAssetsPresen プレゼンテーションのインタフェース
type CryptoAssetsPresen interface {
	Get(w http.ResponseWriter, r *http.Request) error
	// Post(w http.ResponseWriter, r *http.Request) error
	// Put(w http.ResponseWriter, r *http.Request) error
	// Delete(w http.ResponseWriter, r *http.Request) error
}

// cryptoAssetsPresen プレゼンテーションの実体
type cryptoAssetsPresen struct {
	cryptoAssetsUsecase usecase.CryptoAssetsUsecase
}

// CreateCrypocryptoAssetsPresen プレゼンテーションのコンストラクタ
// 戻り値がインタフェースだから 実装を強制できる
func CreateCrypocryptoAssetsPresen(uca usecase.CryptoAssetsUsecase) CryptoAssetsPresen {
	return &cryptoAssetsPresen{uca}
}

func (ca *cryptoAssetsPresen) Get(w http.ResponseWriter, r *http.Request) error {
	// URL パラメータ を string で取得
	paramSymbol := r.URL.Query().Get("symbol")

	averagePriceRate, err := ca.cryptoAssetsUsecase.CoinAveragePrice(paramSymbol)
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
