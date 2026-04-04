package domain

// PasswordHasherRepo は平文パスワードのハッシュ化と照合を担う
type PasswordHasherRepo interface {
	Hash(plain string) (string, error)
	Verify(password Password, plain string) error
}
