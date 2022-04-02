package domain

// TransactionRepository トランザクションのインタフェース
// この struct は ユースケース層 で使われて 実装は インフラ層
type TransactionRepository interface {
	// 作成して どんな値になったか確認できるように 作成した値を戻り値に含める
	// Create(t *Transaction) (*Transaction, error)
	FindByID(transactionId int) (*Transaction, error)
	FindBySymbol(symbol string) (*[]Transaction, error)
	// アップデートして どんな値になったか確認できるように アップデートした値を戻り値に含める
	// Save(t *Transaction) (*Transaction, error)
	// Delete(t *Transaction) error
}
