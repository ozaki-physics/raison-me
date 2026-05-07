package infra

import (
	"fmt"

	"github.com/ozaki-physics/raison-me/share/errorer"
)

// インフラ層の 独自エラー
type InfraError interface {
	error
	Unwrap() error
	FullError() string
}

// 実体
type infraError struct {
	msg string
	err error
}

// エラーを生成する
// ただし エラー自体を生成するとき と ラップして生成するとき があるので プライベートメソッドにした
func newInfraError(msg string, innerErr error) InfraError {
	if msg == "" {
		msg = "インフラ層のエラーを生成しましたが エラーメッセージが格納されていません"
	}
	return &infraError{msg, innerErr}
}

// エラーオブジェクト自体を生成
func NewInfraError(msg string) InfraError {
	return newInfraError(msg, nil)
}

// エラーオブジェクトを インフラ層エラー でラップする
// つまり エラーオブジェクトを生成 と本質は同じ
func WrapInfraError(msg string, innerErr error) InfraError {
	return newInfraError(msg, innerErr)
}

// 標準エラーのインタフェースを満たすため
func (e *infraError) Error() string {
	return e.msg
}

// ラップ元のエラーまで出力する
func (e *infraError) FullError() string {
	// 再帰的に FullError を 呼び出すことで ラップ元のエラーも 全て 出力する
	if e.err != nil {
		if fullErr, ok := e.err.(errorer.FullErrorer); ok {
			return fmt.Sprintf("%s: %s", e.msg, fullErr.FullError())
		}
		return fmt.Sprintf("%s: %v", e.msg, e.err)
	}
	return e.msg
}

// Unwrap したときに ラップ元の型を取り出せるようにするため
func (e *infraError) Unwrap() error {
	return e.err
}
