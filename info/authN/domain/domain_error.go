package domain

import "fmt"

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
func (de *domainError) Error() string {
	return de.msg
}

// ラップ元のエラーまで出力する
func (de *domainError) FullError() string {
	return fmt.Sprintf("%s: %v", de.msg, de.err)
}

// Unwrap したときに ラップ元の型を取り出せるようにするため
func (de *domainError) Unwrap() error {
	return de.err
}
