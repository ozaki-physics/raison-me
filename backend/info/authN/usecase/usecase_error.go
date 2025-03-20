package usecase

import "fmt"

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
func (de *usecaseError) Error() string {
	return de.msg
}

// ラップ元のエラーまで出力する
func (de *usecaseError) FullError() string {
	return fmt.Sprintf("%s: %v", de.msg, de.err)
}

// Unwrap したときに ラップ元の型を取り出せるようにするため
func (de *usecaseError) Unwrap() error {
	return de.err
}
