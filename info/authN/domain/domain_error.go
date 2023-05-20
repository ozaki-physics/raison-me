package domain

// ドメイン層の 独自エラー
type DomainError struct {
	description string
}

func NewDomainError(description string) *DomainError {
	return &DomainError{description}
}

// 標準エラーのインタフェースを満たすため
func (de *DomainError) Error() string {
	return de.description
}
