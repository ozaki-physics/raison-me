package usecase

import (
	"fmt"

	"github.com/ozaki-physics/raison-me/share"
)

// ユースケース層の 独自エラー
type UsecaseError interface {
	error
	Unwrap() error
	FullError() string
}

// 実体
type usecaseError struct {
	msg string
	err error
}

// エラーを生成する
// ただし エラー自体を生成するとき と ラップして生成するとき があるので プライベートメソッドにした
func newUsecaseError(msg string, innerErr error) UsecaseError {
	if msg == "" {
		msg = "ユースケース層のエラーを生成しましたが エラーメッセージが格納されていません"
	}
	return &usecaseError{msg, innerErr}
}

// エラーオブジェクト自体を生成
func NewUsecaseError(msg string) UsecaseError {
	return newUsecaseError(msg, nil)
}

// エラーオブジェクトを ユースケース層エラー でラップする
// つまり エラーオブジェクトを生成 と本質は同じ
func WrapUsecaseError(msg string, innerErr error) UsecaseError {
	return newUsecaseError(msg, innerErr)
}

// 標準エラーのインタフェースを満たすため
func (e *usecaseError) Error() string {
	return e.msg
}

// ラップ元のエラーまで出力する
func (e *usecaseError) FullError() string {
	// 再帰的に FullError を 呼び出すことで ラップ元のエラーも 全て 出力する
	if e.err != nil {
		if fullErr, ok := e.err.(share.FullErrorer); ok {
			return fmt.Sprintf("%s: %s", e.msg, fullErr.FullError())
		}
		return fmt.Sprintf("%s: %v", e.msg, e.err)
	}
	return e.msg
}

// Unwrap したときに ラップ元の型を取り出せるようにするため
func (e *usecaseError) Unwrap() error {
	return e.err
}
