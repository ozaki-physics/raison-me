package infrastructure

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/ozaki-physics/raison-me/capital/domain"
)

// data JSON から struct に変換する
type dataCMCID struct {
	CMCIDs []struct {
		Symbol string `json:"symbol"`
		CMCID  int    `json:"CMCID"`
	} `json:"data"`
}

// coinInfra コインインフラの実体
type coinInfra struct {
	data dataCMCID
}

// CreateCoinRepository コインインフラのコンストラクタ
// 戻り値がインタフェースだから 実装を強制できる
func CreateCoinRepository() domain.CoinRepository {
	// 何度も JSON を読み込まなくていいように インスタンス変数に格納しておく

	// CoinMarketCap ID の対応表を読み込む
	bytes, err := os.ReadFile("./capital/json/cryptoassets_Coin.json")
	if err != nil {
		log.Fatalln(err)
	}
	// JSON から struct にする
	var d dataCMCID
	if err := json.Unmarshal(bytes, &d); err != nil {
		log.Fatalln(err)
	}

	return &coinInfra{d}
}

// func (c *CoinInfra) Create(dc *domain.Coin) (*domain.Coin, error) {}

func (c *coinInfra) FindBySymbol(symbol string) (*domain.Coin, error) {
	for _, cmcIdtable := range c.data.CMCIDs {
		if cmcIdtable.Symbol == symbol {

			// CoinMarketCap にリクエストして 現在価格を取得する
			var CMCIDs []int
			CMCIDs = append(CMCIDs, cmcIdtable.CMCID)
			cmc := getCredential()
			// 1個しか引数に渡してないから 1個しか返却されてないと仮定する
			price := cmc.getQuotesLatest(CMCIDs)[0].Price

			dc, err := domain.ReconstructCoin(symbol, price)
			if err != nil {
				return nil, err
			}
			return dc, nil
		}
	}
	return nil, errors.New("CoinMarketCap ID が存在しません")
}

// func (c *CoinInfra) Save(dc *domain.Coin) (*domain.Coin, error) {}

// func (c *CoinInfra) Delete(dc *domain.Coin) error {}
