package domain

import "context"

// この struct は ユースケース層 で使われていて 実装は インフラ層
type PassRepo interface {
	passRepoInsert
	passRepoFetch
	// passRepoUpdate
	// passRepoDelete
}

// 追加する
type passRepoInsert interface {
	Insert(ctx context.Context, pass Pass) (*Pass, error)
}

// 取得する
type passRepoFetch interface {
	FindByAccountId(ctx context.Context, accountID AccountID) (*Pass, error)
}

// 更新する
type passRepoUpdate interface {
	Update(ctx context.Context, pass Pass) (*Pass, error)
}

// 削除する
type passRepoDelete interface {
	Delete(ctx context.Context, accountID AccountID) error
}
