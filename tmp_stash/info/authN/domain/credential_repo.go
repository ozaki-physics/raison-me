package domain

// この struct は ユースケース層 で使われていて 実装は インフラ層
type CredentialRepo interface {
	credentialRepoInsert
	credentialRepoFetch
	// credentialRepoUpdate
	// credentialRepoDelete
}

// 追加する
type credentialRepoInsert interface {
	Insert(c Credential) (*Credential, DomainError)
}

// 取得する
type credentialRepoFetch interface {
	Fetch() ([]Credential, DomainError)
	FindByAccountId(accountID AccountID) (*Credential, DomainError)
	// FindById(id UserID) (*Credential, DomainError)
	// FindByName(name UserName) (*Credential, DomainError)
}

// 更新する
type credentialRepoUpdate interface {
	Update(c *Credential) (*Credential, DomainError)
}

// 削除する
type credentialRepoDelete interface {
	Delete(accountID AccountID) DomainError
}
