package presen

import (
	"fmt"

	"github.com/ozaki-physics/raison-me/share/errorer"
)

// プレゼンテーション層の 独自エラー
type PresenError interface {
	error
	Unwrap() error
	FullError() string
	// ステータスコードを 指定するため
	StatusCode() int
}

// 実体
type presenError struct {
	msg        string
	err        error
	statusCode int
}

// エラーを生成する
// ただし エラー自体を生成するとき と ラップして生成するとき があるので プライベートメソッドにした
func newPresenError(msg string, innerErr error, statusCode int) PresenError {
	if msg == "" {
		msg = "プレゼンテーション層のエラーを生成しましたが エラーメッセージが格納されていません"
	}
	return &presenError{msg, innerErr, statusCode}
}

// エラーオブジェクト自体を生成
func NewPresenError(msg string, statusCode int) PresenError {
	return newPresenError(msg, nil, statusCode)
}

// エラーオブジェクトを プレゼンテーション層エラー でラップする
// つまり エラーオブジェクトを生成 と本質は同じ
func WrapPresenError(msg string, innerErr error, statusCode int) PresenError {
	return newPresenError(msg, innerErr, statusCode)
}

// 標準エラーのインタフェースを満たすため
func (e *presenError) Error() string {
	return e.msg
}

// ラップ元のエラーまで出力する
func (e *presenError) FullError() string {
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
func (e *presenError) Unwrap() error {
	return e.err
}

// ステータスコードを返す
func (e *presenError) StatusCode() int {
	return e.statusCode
}
