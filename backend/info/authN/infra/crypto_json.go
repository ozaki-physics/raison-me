package infra

import (
	"github.com/ozaki-physics/raison-me/info/authN/domain"
	"github.com/ozaki-physics/raison-me/share"
)

const cryptoPath = "./info/authN/infra/json/salt.json"

type cryptoRepoJSON struct {
	// 何度も JSON を読み込まなくていいように インスタンス変数に格納しておく?
	No01 string
	No02 string
}

func NewCryptoRepoJSON() (domain.CryptoRepo, error) {
	var jc jsonCrypto
	// main.go からの path
	share.JsonToStruct(&jc, cryptoPath)

	// JSON を インスタンス に保持させるため
	c := &cryptoRepoJSON{
		No01: jc.Salt.No01,
		No02: jc.Salt.No02,
	}
	return c, nil
}

func (c *cryptoRepoJSON) Read() (*domain.Crypto, error) {
	dc, err := domain.NewCrypto(c.No01, c.No02)
	if err != nil {
		return nil, err
	}
	return dc, nil
}

type jsonCrypto struct {
	Salt salt `json:"salt"`
}

type salt struct {
	No01 string `json:"no01"`
	No02 string `json:"no02"`
}
