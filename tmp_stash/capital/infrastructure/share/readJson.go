package share

import (
	"encoding/json"
	"log"
	"os"
)

func JsonToStruct(goStruct interface{}, path string) {
	// apiKey を読み込む
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	// JSON から struct にする
	if err := json.Unmarshal(bytes, &goStruct); err != nil {
		log.Fatalln(err)
	}
}
