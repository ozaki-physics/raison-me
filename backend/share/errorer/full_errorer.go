package errorer

// FullErrorer を 実装している エラー なら FullError を 返すため
type FullErrorer interface {
	FullError() string
}
