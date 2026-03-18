package domain

import "context"

// この struct は ユースケース層 で使われていて 実装は インフラ層
type UserRepo interface {
	userRepoInsert
	userRepoFetch
	// userRepoUpdate
	// userRepoDelete
}

// 追加する
type userRepoInsert interface {
	Insert(ctx context.Context, user *User) (*User, error)
}

// 取得する
type userRepoFetch interface {
	Fetch(ctx context.Context) ([]User, error)
	FindByAccountId(ctx context.Context, accountID AccountID) (*User, error)
	FindById(ctx context.Context, id UserID) (*User, error)
	// FindByName(ctx context.Context, name UserName) (*User, error)
}

// 更新する
type userRepoUpdate interface {
	Update(ctx context.Context, user *User) (*User, error)
}

// 削除する
type userRepoDelete interface {
	Delete(ctx context.Context, accountID AccountID) error
}
