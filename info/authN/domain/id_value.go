package domain

import (
	"strings"

	"github.com/rs/xid"
)

type id struct {
	prefixID string
	value    string
}

func constructorID(prefixID string, val string) (id, error) {
	if prefixID == "" {
		return NilID(), NewDomainError("ID生成時のプレフィックスが存在しません")
	}
	if val == "" {
		return NilID(), NewDomainError("ID生成時の値が存在しません")
	}

	id := id{prefixID, val}
	if len(id.Val()) > 30 {
		return id, NewDomainError("ID生成時のプレフィックスが長すぎます")
	}

	return id, nil
}

func NewID(prefixID string) (id, error) {
	guid := xid.New()
	return constructorID(prefixID, guid.String())
}

func ReNewID(data string) (id, error) {
	sp := strings.Split(data, "-")
	var val string
	if len(sp) >= 3 {
		val = strings.Join(sp[1:], "-")
	} else {
		val = sp[1]
	}
	return constructorID(sp[0], val)
}

func NilID() id {
	nilID, _ := constructorID("nil", "nil")
	return nilID
}

// 以下ゲッター
func (id *id) Val() string {
	return id.prefixID + "-" + id.value
}

func (id *id) PrefixID() string {
	return id.prefixID
}
