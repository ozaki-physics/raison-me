package domain

import (
	"crypto/sha512"
	"hash"
	"io"
)

type Crypto struct {
	salt01        string
	salt02        string
	hashAlgorithm hash.Hash
}

func NewCrypto(salt01 string, salt02 string) (*Crypto, error) {
	// TODO: err が未実装
	c := &Crypto{
		salt01:        salt01,
		salt02:        salt02,
		hashAlgorithm: sha512.New(),
	}
	return c, nil
}

func (c *Crypto) Hash(name UserName, pass Password) (string, error) {
	// TODO: io.WriteString() の 戻り値を無視してるけどいいのかな
	// TODO: 本当にこのアルゴリズムでいいのか?
	// TODO: 長さとか他に気をつける必要があることはなんだろう?
	// TODO: salt の文字列はユーザーごとに変えた方がいいらしいけど どうやって salt を保持するの?
	io.WriteString(c.hashAlgorithm, pass.ToHash())

	hash01Byte := c.hashAlgorithm.Sum(nil)
	hash01String := string(hash01Byte)

	io.WriteString(c.hashAlgorithm, c.salt01)
	io.WriteString(c.hashAlgorithm, name.Val())
	io.WriteString(c.hashAlgorithm, c.salt02)
	io.WriteString(c.hashAlgorithm, hash01String)
	hash02Byte := c.hashAlgorithm.Sum(nil)
	return string(hash02Byte), nil
}
