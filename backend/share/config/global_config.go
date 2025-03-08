package config

import "os"

type Config interface {
	IsLive() bool
	IsCloud() bool
	GetGCPProjectID() string
}

const (
	ProductionCloud = iota
	ProductionLocal
	DevelopmentCloud
	DevelopmentLocal
)

const runMode = DevelopmentCloud

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
	runMode      int
	isLive       bool
	isCloud      bool
	gcpProjectID string
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

func (c *config) GetGCPProjectID() string {
	tmp := os.Getenv("GCP_PROJECT_ID")
	if tmp != "" {
		return tmp
	}
	return c.gcpProjectID
}

// 本番(実データ, クラウド)
func newProductionCloudConfig() *config {
	c := &config{
		runMode:      ProductionCloud,
		isLive:       true,
		isCloud:      true,
		gcpProjectID: "raison-me",
	}
	return c
}

// 本番(実データ, ローカル)
func newProductionLocalConfig() *config {
	c := &config{
		runMode:      ProductionLocal,
		isLive:       true,
		isCloud:      false,
		gcpProjectID: "",
	}
	return c
}

// 開発(テストデータ, クラウド)
func newDevelopmentCloudConfig() *config {
	c := &config{
		runMode:      DevelopmentCloud,
		isLive:       false,
		isCloud:      true,
		gcpProjectID: "smart-ruler-277318",
	}
	return c
}

// 開発(テストデータ, ローカル)
func newDevelopmentLocalConfig() *config {
	c := &config{
		runMode:      DevelopmentLocal,
		isLive:       false,
		isCloud:      false,
		gcpProjectID: "",
	}
	return c
}
