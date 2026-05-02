package presen

import (
	"fmt"
)

// プレゼンテーション層の 独自エラー
type PresenError interface {
	error
	Unwrap() error
	FullError() string
	// ステータスコードを 指定するため
	StatusCoder() int
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
func (de *presenError) Error() string {
	return de.msg
}

// ラップ元のエラーまで出力する
func (de *presenError) FullError() string {
	return fmt.Sprintf("%s: %v", de.msg, de.err)
}

// Unwrap したときに ラップ元の型を取り出せるようにするため
func (de *presenError) Unwrap() error {
	return de.err
}

// ステータスコードを返す
func (de *presenError) StatusCoder() int {
	return de.statusCode
}
