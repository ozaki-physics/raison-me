package infra

import (
	"encoding/json"
	"log"
	"os"

	c "github.com/ozaki-physics/raison-me/share/config"
)

// ストレージ(JSON)などの main.go からの path を保持
type storagePath struct {
	userJSON string
	passJSON string
}

// 環境によってパス先を変えるため
func NewStoragePath() storagePath {
	globalConfig := c.NewConfig()

	if globalConfig.IsCloud() {
		if globalConfig.IsLive() {
			return storagePath{
				userJSON: "./persist/user_persist_example.json",
				passJSON: "./persist/pass_persist_example.json",
			}
		} else {
			return storagePath{
				userJSON: "./persist/user_example.json",
				passJSON: "./persist/pass_example.json",
			}
		}
	} else {
		if globalConfig.IsLive() {
			return storagePath{
				userJSON: "info/authN/infra/json/user_persist_example.json",
				passJSON: "info/authN/infra/json/pass_persist_example.json",
			}
		} else {
			return storagePath{
				userJSON: "info/authN/infra/json/user_example.json",
				passJSON: "info/authN/infra/json/pass_example.json",
			}
		}
	}
}

// TODO: マジで暫定的 な ファイル区別のためのユーザー取得
func (s *storagePath) GetUser() int {
	data, err := os.ReadFile(s.userJSON)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	var users struct {
		User []struct {
			AccountID string `json:"account_id"`
			ID        string `json:"id"`
			Name      string `json:"name"`
		} `json:"user"`
	}
	if err := json.Unmarshal(data, &users); err != nil {
		log.Fatalf("failed to unmarshal JSON: %v", err)
	}

	return len(users.User)
}

// TODO: 用途不明
func (s *storagePath) GetUserPath() string {
	return s.userJSON
}
