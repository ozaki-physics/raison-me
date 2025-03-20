package domain

import (
	"strings"

	"github.com/rs/xid"
)

type id struct {
	prefixID string
	value    string
}

func constructorID(prefixID string, val string) (id, DomainError) {
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

func NewID(prefixID string) (id, DomainError) {
	guid := xid.New()
	return constructorID(prefixID, guid.String())
}

func ReNewID(data string) (id, DomainError) {
	sp := strings.Split(data, "-")
	var val string

	// 本来は 値オブジェクトが作れた時点で不整合な状態にはならないから
	// 復元するときにも不整合な状態は発生しないはず
	// もしも ストレージに直接 ID を書き込んで それが不整合だったときを考える
	if len(sp) == 1 {
		return NilID(), NewDomainError("ストレージのIDフォーマットが不正です")
	}

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

func (id *id) IsNilID() bool {
	return *id == NilID()
}

// 以下ゲッター

func (id *id) Val() string {
	return id.prefixID + "-" + id.value
}

func (id *id) PrefixID() string {
	return id.prefixID
}
