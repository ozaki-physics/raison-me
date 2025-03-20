package domain

// この struct は ユースケース層 で使われていて 実装は インフラ層
type UserRepo interface {
	userRepoInsert
	userRepoFetch
	// userRepoUpdate
	// userRepoDelete
}

// 追加する
type userRepoInsert interface {
	Insert(id UserID, name UserName) (*User, error)
}

// 取得する
type userRepoFetch interface {
	Fetch() ([]User, error)
	FindByAccountId(accountID AccountID) (*User, error)
	FindById(id UserID) (*User, error)
	FindByName(name UserName) (*User, error)
}

// 更新する
type userRepoUpdate interface {
	Update(user *User) (*User, error)
}

// 削除する
type userRepoDelete interface {
	Delete(accountID AccountID) error
}
