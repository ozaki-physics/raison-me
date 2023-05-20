package domain

import (
	"github.com/rs/xid"
)

type id struct {
	prefixID string
	value    xid.ID
}

func NewID(prefixID string) (id, error) {
	// バリデーションが先だが 値オブジェクトを返すために id が nil だとダメだから
	guid := xid.New()
	id := id{prefixID, guid}

	if prefixID == "" {
		return id, NewDomainError("ID生成時のプレフィックスが存在しません")
	}

	if len(id.Val()) > 30 {
		return id, NewDomainError("ID生成時のプレフィックスが長すぎます")
	}

	return id, nil
}

// 以下ゲッター
func (id *id) Val() string {
	return id.prefixID + "-" + id.value.String()
}

func (id *id) PrefixID() string {
	return id.prefixID
}
