package infra

import "fmt"

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
func (ie *infraError) Error() string {
	return ie.msg
}

// ラップ元のエラーまで出力する
func (ie *infraError) FullError() string {
	return fmt.Sprintf("%s: %v", ie.msg, ie.err)
}

// Unwrap したときに ラップ元の型を取り出せるようにするため
func (ie *infraError) Unwrap() error {
	return ie.err
}
