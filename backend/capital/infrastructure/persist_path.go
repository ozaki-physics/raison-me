package infrastructure

import (
	c "github.com/ozaki-physics/raison-me/share/config"
)

type Persist_path interface {
	GetCoinMarketCapId() string
	GetTransaction() string
}

func NewConfig() Persist_path {
	globalConfig := c.NewConfig()

	if globalConfig.IsCloud() {
		return &config{
			coinMarketCapId: "./persist/coinMarketCapId.json",
			transaction:     "./persist/transaction.json",
		}
	} else {
		return &config{
			coinMarketCapId: "./capital/infrastructure/cryptoAsset/json/coinMarketCapId.json",
			transaction:     "./capital/infrastructure/cryptoAsset/json/transaction.json",
		}
	}
}

type config struct {
	coinMarketCapId string
	transaction     string
}

func (c *config) GetCoinMarketCapId() string {
	return c.coinMarketCapId
}

func (c *config) GetTransaction() string {
	return c.transaction
}
