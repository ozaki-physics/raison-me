package domain

// API キー 値オブジェクト
type APIKey string

func NewAPIKey(key string) (APIKey, error) {
	// TODO: 値としてのバリデーションが未実装
	// 使える文字種の制限
	k := APIKey(key)
	return k, nil
}

func (k *APIKey) toString() string {
	return string(*k)
}

// 機密フィールドの出力形式変更
func (k *APIKey) String() string {
	return "xxxxxx"
}

// 機密フィールドの出力形式変更
func (k *APIKey) GoString() string {
	return "xxxxxx"
}
