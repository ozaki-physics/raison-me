package cryptoasset

import (
	"errors"
	"math"
)

type Transaction struct {
	symbol    string  // 銘柄
	side      int     // 売買区分 0:買 1:売
	priceRate float64 // 約定レート
	size      float64 // 約定数量
	fee       int     // 取引手数料
	time      string  // 約定日時
}

// constructTransaction コンストラクタ
func constructTransaction(symbol string, side int, priceRate float64, size float64, fee int, time string) (*Transaction, error) {
	if symbol == "" {
		return nil, errors.New("銘柄を登録してください")
	}
	if side != 0 && side != 1 {
		return nil, errors.New("売買区分を登録してください")
	}
	if priceRate == 0 {
		return nil, errors.New("約定レートを登録してください")
	}
	if size == 0 {
		return nil, errors.New("約定数量を登録してください")
	}

	t := &Transaction{
		symbol,
		side,
		priceRate,
		size,
		fee,
		time,
	}
	return t, nil
}

// createTransaction 取引履歴を作成
// func createTransaction(symbol string, side int, priceRate float64, size float64, fee int, time string) (*Transaction, error) {
// 	return constructTransaction(symbol, side, priceRate, size, fee, time)
// }

// reconstructTransaction DB などの値からインスタンスを再構成
func ReconstructTransaction(symbol string, side int, priceRate float64, size float64, fee int, time string) (*Transaction, error) {
	return constructTransaction(symbol, side, priceRate, size, fee, time)
}

// TransactionPrice 約定代金を取得
func (t *Transaction) TransactionPrice() float64 {
	price := t.priceRate * float64(t.size)
	return math.Ceil(price)
}

// ゲッター
func (t *Transaction) PriceRate() float64 {
	return t.priceRate
}
func (t *Transaction) Size() float64 {
	return t.size
}
