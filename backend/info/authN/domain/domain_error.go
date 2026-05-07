package domain

import (
	"fmt"

	"github.com/ozaki-physics/raison-me/share/errorer"
)

// ドメイン層の 独自エラー
type DomainError interface {
	error
	Unwrap() error
	FullError() string
}

// 実体
type domainError struct {
	msg string
	err error
}

// エラーを生成する
// ただし エラー自体を生成するとき と ラップして生成するとき があるので プライベートメソッドにした
func newDomainError(msg string, innerErr error) DomainError {
	if msg == "" {
		msg = "ドメイン層のエラーを生成しましたが エラーメッセージが格納されていません"
	}
	return &domainError{msg, innerErr}
}

// エラーオブジェクト自体を生成
func NewDomainError(msg string) DomainError {
	return newDomainError(msg, nil)
}

// エラーオブジェクトを ドメイン層エラー でラップする
// つまり エラーオブジェクトを生成 と本質は同じ
func WrapDomainError(msg string, innerErr error) DomainError {
	return newDomainError(msg, innerErr)
}

// 標準エラーのインタフェースを満たすため
func (e *domainError) Error() string {
	return e.msg
}

// ラップ元のエラーまで出力する
func (e *domainError) FullError() string {
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
func (e *domainError) Unwrap() error {
	return e.err
}
