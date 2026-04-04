package domain

import "github.com/google/uuid"

type id struct {
	value uuid.UUID
}

func constructorID(val uuid.UUID) (id, DomainError) {
	if val == uuid.Nil {
		return NilID(), NewDomainError("ID生成時の値が存在しません")
	}
	if val.Version() != 7 {
		return NilID(), NewDomainError("ID は UUID v7 の必要があります")
	}

	return id{value: val}, nil
}

func NewID() (id, DomainError) {
	guid, err := uuid.NewV7()
	if err != nil {
		return NilID(), NewDomainError("ID の生成に失敗しました")
	}

	return constructorID(guid)
}

func ReNewID(data string) (id, DomainError) {
	if data == "" {
		return NilID(), NewDomainError("ID 生成時の値が存在しません")
	}

	// 本来は 値オブジェクトが作れた時点で不整合な状態にはならないから
	// 復元するときにも不整合な状態は発生しないはず
	// もしも ストレージに直接 ID を書き込んで それが不整合だったときを考える
	guid, err := uuid.Parse(data)
	if err != nil {
		return NilID(), NewDomainError("ストレージの ID フォーマットが不正です")
	}

	return constructorID(guid)
}

func NilID() id {
	return id{value: uuid.Nil}
}

func (id *id) IsNilID() bool {
	return id.value == uuid.Nil
}

// 以下ゲッター

func (id *id) Val() string {
	return id.value.String()
}
