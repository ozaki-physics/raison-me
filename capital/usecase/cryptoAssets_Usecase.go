package usecase

import (
	"log"

	"github.com/ozaki-physics/raison-me/capital/domain"
)

// CryptoAssetsUsecase ユースケースのインタフェース
type CryptoAssetsUsecase interface {
	// 平均約定レート を確認する
	CoinAveragePrice(symbol string) (float64, error)
	// 評価額 を確認する
	CoinPrice(symbol string) (float64, error)
	// 保有枚数 を確認する
	CoinSize(symbol string) (float64, error)
	// コインの損益 を確認する
	CoinGainPrice(symbol string) (float64, error)
	// 損益の割合 を確認する
	CoinGainPercent(symbol string) (float64, error)
	// 価格ごとの枚数 を確認する
	CoinPriceStepSize(symbol string) (map[float64]float64, error)
	// ある取引の約定代金 を確認する
	TransactionPrice(transactionId int) (float64, error)
}

// cryptoAssetsUsecase ユースケースの実体
type cryptoAssetsUsecase struct {
	coin        domain.CoinRepository
	transaction domain.TransactionRepository
}

// CreateCryptoAssetsUsecase ユースケースのコンストラクタ
// 戻り値がインタフェースだから 実装を強制できる
func CreateCryptoAssetsUsecase(dc domain.CoinRepository, dt domain.TransactionRepository) CryptoAssetsUsecase {
	return &cryptoAssetsUsecase{dc, dt}
}

// 平均約定レート を確認する
func (u *cryptoAssetsUsecase) CoinAveragePrice(symbol string) (float64, error) {
	manyTransaction, err := u.transaction.FindBySymbol(symbol)
	if err != nil {
		return 0, err
	}

	var averagePriceRate float64
	var sumSize float64
	for _, t := range *manyTransaction {
		averagePriceRate = (averagePriceRate*sumSize + t.PriceRate()*t.Size()) / (sumSize + t.Size())
		sumSize += t.Size()
	}
	return averagePriceRate, nil
}

// 評価額 を確認する
func (u *cryptoAssetsUsecase) CoinPrice(symbol string) (float64, error) {
	averagePrice, err := u.CoinAveragePrice(symbol)
	if err != nil {
		return 0, err
	}
	sumSize, err := u.CoinSize(symbol)
	if err != nil {
		return 0, err
	}
	return averagePrice * sumSize, nil
}

// 保有枚数 を確認する
func (u *cryptoAssetsUsecase) CoinSize(symbol string) (float64, error) {
	manyTransaction, err := u.transaction.FindBySymbol(symbol)
	if err != nil {
		log.Fatal(err)
	}

	var sumSize float64
	for _, t := range *manyTransaction {
		sumSize += t.Size()
	}
	return sumSize, nil
}

// コインの損益 を確認する
func (u *cryptoAssetsUsecase) CoinGainPrice(symbol string) (float64, error) {
	nowCoin, err := u.coin.FindBySymbol(symbol)
	if err != nil {
		return 0, err
	}
	averagePrice, err := u.CoinAveragePrice(symbol)
	if err != nil {
		return 0, err
	}
	size, err := u.CoinSize(symbol)
	if err != nil {
		return 0, err
	}
	diff := nowCoin.Bid() - averagePrice
	return diff * size, nil
}

// 損益の割合 を確認する
func (u *cryptoAssetsUsecase) CoinGainPercent(symbol string) (float64, error) {
	gainPrice, err := u.CoinGainPrice(symbol)
	if err != nil {
		return 0, err
	}
	price, err := u.CoinPrice(symbol)
	if err != nil {
		return 0, err
	}
	return gainPrice / price, nil
}

// 価格ごとの枚数 を確認する
func (u *cryptoAssetsUsecase) CoinPriceStepSize(symbol string) (map[float64]float64, error) {
	manyTransaction, err := u.transaction.FindBySymbol(symbol)
	if err != nil {
		return nil, err
	}
	var step map[float64]float64
	for _, t := range *manyTransaction {
		step[t.PriceRate()] += t.Size()
	}
	return step, nil
}

// ある取引の約定代金 を確認する
func (u *cryptoAssetsUsecase) TransactionPrice(transactionId int) (float64, error) {
	t, err := u.transaction.FindByID(transactionId)
	if err != nil {
		return 0, err
	}
	return t.TransactionPrice(), nil
}
