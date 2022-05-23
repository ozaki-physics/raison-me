package cryptoasset

import "errors"

type Coin struct {
	symbol string  // 銘柄
	bid    float64 // 現在値(売値)
}

// constructCoin コンストラクタ
func constructCoin(symbol string, bid float64) (*Coin, error) {
	if symbol == "" {
		return nil, errors.New("銘柄を登録してください")
	}

	c := &Coin{
		symbol,
		bid,
	}
	return c, nil
}

// func createCoin(symbol string) (*Coin, error) {
// return constructCoin(symbol)
// }

// reconstructCoin DB などの値からインスタンスを再構成
func ReconstructCoin(symbol string, bid float64) (*Coin, error) {
	return constructCoin(symbol, bid)
}

// ゲッター
func (c *Coin) Bid() float64 {
	return c.bid
}
