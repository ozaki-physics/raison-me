package cryptoasset

import (
	"encoding/json"
	"log"
	"os"

	infra "github.com/ozaki-physics/raison-me/capital/infrastructure"
)

// dataCMCIds JSON から struct に変換する
type dataCMCIds struct {
	Data []struct {
		Symbol string `json:"symbol"`
		CMCId  int    `json:"CMCId"`
	} `json:"data"`
}

func CreateCMCIdsJson() CMCIds {
	// 何度も JSON を読み込まなくていいように インスタンス変数に格納しておく

	persist := infra.NewConfig()
	// apiKey を読み込む
	bytes, err := os.ReadFile(persist.GetCoinMarketCapId())
	if err != nil {
		log.Fatalln(err)
	}
	// JSON から struct にする
	var d dataCMCIds
	if err := json.Unmarshal(bytes, &d); err != nil {
		log.Fatalln(err)
	}

	var cmdIdMap = make(map[string]int)
	for _, cmcMap := range d.Data {
		cmdIdMap[cmcMap.Symbol] = cmcMap.CMCId
	}

	return &coinMarketCapIdDto{cmdIdMap}
}

func (c *coinMarketCapIdDto) SymbolAndId() map[string]int {
	return c.cmcIdMap
}

func (c *coinMarketCapIdDto) Id(symbol string) int {
	return c.cmcIdMap[symbol]
}
