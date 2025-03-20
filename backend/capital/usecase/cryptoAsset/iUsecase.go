package cryptoasset

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
