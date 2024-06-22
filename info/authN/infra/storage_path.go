package infra

import "github.com/ozaki-physics/raison-me/share/config"

// ストレージ(JSON)などの main.go からの path を保持
type storagePath struct {
	userJSON string
	passJSON string
}

// 環境によってパス先を変えるため
func newStoragePath() storagePath {
	s := storagePath{}

	if config.IsCloud == false {
		if config.IsLive {
			s.userJSON = "info/authN/infra/json/user.json"
			s.passJSON = "info/authN/infra/json/pass.json"
			return s
		} else {
			s.userJSON = "info/authN/infra/json/user_sample.json"
			s.passJSON = "info/authN/infra/json/pass_sample.json"
		}
	}
	return s
}
