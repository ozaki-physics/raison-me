package cryptoasset

// CoinRepository コインのインタフェース
// この struct は ユースケース層 で使われて 実装は インフラ層
type CoinRepository interface {
	// Create(c *Coin) (*Coin, error)
	FindBySymbol(symbol string) (*Coin, error)
	// Save(c *Coin) (*Coin, error)
	// Delete(c *Coin) error
}
