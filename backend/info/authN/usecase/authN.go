package usecase

import "github.com/ozaki-physics/raison-me/info/authN/domain"

// 認証コンテキスト で 必要なユースケース
type AuthN interface {
	User
}

// 認証コンテキスト で 必要なユースケース(User 特化)
type User interface {
	// ユーザー を 作成
	Create(userDto UserDto) (*UserDto, UsecaseError)
	// ユーザー を 削除
	// Delete(userDto UserDto) (bool, UsecaseError)
	// ユーザー で ログイン
	SignIn(userDto UserDto) (bool, UsecaseError)
	// ユーザー で ログアウト
	SignOut(userDto UserDto) (bool, UsecaseError)
	// ユーザー(自分) の情報を更新
	// Update(userDto UserDto) (*UserDto, UsecaseError)
	// パスワードリセット
	// ResetPassword()
}

type authN struct {
	credentialRepo domain.CredentialRepo
}

func NewAuthN(dcr domain.CredentialRepo) (AuthN, UsecaseError) {
	u := &authN{dcr}
	return u, nil
}

func (a *authN) Create(userDto UserDto) (*UserDto, UsecaseError) {
	// User の作成
	aID, err := domain.NewAccountID()
	if err != nil {
		return nil, err
	}
	uID, err := domain.NewUserID(userDto.ID())
	if err != nil {
		return nil, err
	}
	uName, err := domain.NewUserName(userDto.Name())
	if err != nil {
		return nil, err
	}
	u, err := domain.NewUser(aID, uID, uName)
	if err != nil {
		return nil, err
	}

	// Pass の作成
	pID, err := domain.NewPassID()
	if err != nil {
		return nil, err
	}
	iat, err := domain.NewDate()
	if err != nil {
		return nil, err
	}
	password, err := domain.NewPassword(userDto.password())
	if err != nil {
		return nil, err
	}
	p, err := domain.NewPass(pID, aID, password, iat)
	if err != nil {
		return nil, err
	}

	// Credential の作成
	c, err := domain.NewCredential(*u, *p)
	if err != nil {
		return nil, err
	}
	saveCredential, err02 := a.credentialRepo.Insert(*c)
	if err02 != nil {
		return nil, WrapUsecaseError("保存に失敗しました", err02)
	}

	// プレゼン層に返す値を生成
	saveU := saveCredential.User()
	saveP := saveCredential.Pass()
	saveUserDto := NewUserDto(saveU.AccountID(), saveU.ID(), saveU.Name(), saveP.Password())
	return saveUserDto, nil
}

func (a *authN) SignIn(userDto UserDto) (bool, UsecaseError) {
	aID, err := domain.ReNewAccountID(userDto.AccountID())
	if err != nil {
		return false, WrapUsecaseError("ログインに失敗しました", err)
	}
	c, err := a.credentialRepo.FindByAccountId(aID)
	if err != nil {
		return false, WrapUsecaseError("ログインに失敗しました", err)
	}
	pass := c.Pass()
	isOK, err := pass.IsLogin(userDto.password())
	if err != nil {
		return false, WrapUsecaseError("ログインに失敗しました", err)
	}

	return isOK, nil
}

func (a *authN) SignOut(userDto UserDto) (bool, UsecaseError) {
	// TODO: ログアウトの処理とは?
	isOK := true
	return isOK, nil
}

// TODO: 用途不明
func NewUser(userRepo domain.UserRepo, cryptoRepo domain.CryptoRepo) (User, UsecaseError) {
	return nil, newUsecaseError("未実装", nil)
}
