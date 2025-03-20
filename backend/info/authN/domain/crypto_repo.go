package domain

// この struct は ユースケース層 で使われていて 実装は インフラ層
type CryptoRepo interface {
	// salt の値を取得してくる
	Read() (*Crypto, error)
}
