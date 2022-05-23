package cryptoasset

type CMCIds interface {
	// インタフェースにして 実体の型が失われていて 外部で使うときに型アサーションするのも違う気がしたから(実際には 実体をそのまま返しているだけ)
	SymbolAndId() map[string]int
	Id(symbol string) int
}

// coinMarketCapIdDto CoinMarketCap の id を持つ実体
type coinMarketCapIdDto struct {
	// struct の スライスでもよいが map だと symbol 渡されたら 1発で 探し当てられそうだから
	cmcIdMap map[string]int
}
