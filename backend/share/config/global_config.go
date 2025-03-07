package config

import "os"

type Config interface {
	IsLive() bool
	IsCloud() bool
}

const (
	ProductionCloud = iota
	ProductionLocal
	DevelopmentCloud
	DevelopmentLocal
)

const runMode = DevelopmentLocal

func NewConfig() Config {
	switch runMode {
	case ProductionCloud:
		return newProductionCloudConfig()
	case ProductionLocal:
		return newProductionLocalConfig()
	case DevelopmentCloud:
		return newDevelopmentCloudConfig()
	case DevelopmentLocal:
		return newDevelopmentLocalConfig()
	default:
		return newDevelopmentLocalConfig()
	}
}

type config struct {
	runMode int
	isLive  bool
	isCloud bool
}

func (c *config) IsLive() bool {
	// TODO: 無理やり環境変数から取得している
	if os.Getenv("IS_LIVE") == "true" {
		return true
	}
	return c.isLive
}

func (c *config) IsCloud() bool {
	// TODO: 無理やり環境変数から取得している
	if os.Getenv("IS_CLOUD") == "true" {
		return true
	}
	return c.isCloud
}

// 本番(実データ, クラウド)
func newProductionCloudConfig() *config {
	c := &config{
		runMode: ProductionCloud,
		isLive:  true,
		isCloud: true,
	}
	return c
}

// 本番(実データ, ローカル)
func newProductionLocalConfig() *config {
	c := &config{
		runMode: ProductionLocal,
		isLive:  true,
		isCloud: false,
	}
	return c
}

// 開発(テストデータ, クラウド)
func newDevelopmentCloudConfig() *config {
	c := &config{
		runMode: DevelopmentCloud,
		isLive:  false,
		isCloud: true,
	}
	return c
}

// 開発(テストデータ, ローカル)
func newDevelopmentLocalConfig() *config {
	c := &config{
		runMode: DevelopmentLocal,
		isLive:  false,
		isCloud: false,
	}
	return c
}
