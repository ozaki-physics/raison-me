package usecase

import (
	"context"

	"github.com/ozaki-physics/raison-me/info/authN/domain"
)

// 認証コンテキスト で 必要なユースケース
type AuthN interface {
	User
}

// 認証コンテキスト で 必要なユースケース(User 特化)
type User interface {
	// ユーザー を 作成
	Create(ctx context.Context, userName string, userID string, passPlane string) (*UserDto, UsecaseError)
	// ユーザー を 削除
	// Delete(userDto UserDto) (bool, UsecaseError)
	// ユーザー で サインイン
	SignIn(ctx context.Context, userID string, passwordPlane string) (bool, UsecaseError)
	// ユーザー で サインアウト
	SignOut(ctx context.Context, userID string) (bool, UsecaseError)
	// ユーザー(自分) の情報を更新
	// Update(userDto UserDto) (*UserDto, UsecaseError)
	// パスワードリセット
	// ResetPassword()
	// ユーザー一覧 を 取得
	GetUserList(ctx context.Context) ([]UserDto, UsecaseError)
	// ユーザー を 検索
	SearchUser(ctx context.Context, userID string) (*UserDto, UsecaseError)
}

type authN struct {
	userRepo domain.UserRepo
	passRepo domain.PassRepo
}

func NewAuthN(u domain.UserRepo, p domain.PassRepo) (AuthN, UsecaseError) {
	a := &authN{u, p}
	return a, nil
}

func (a *authN) Create(ctx context.Context, userName string, userID string, passPlane string) (*UserDto, UsecaseError) {
	// User の作成
	aID, err := domain.NewAccountID()
	if err != nil {
		return nil, err
	}
	uID, err := domain.NewUserID(userID)
	if err != nil {
		return nil, err
	}
	uName, err := domain.NewUserName(userName)
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
	password, err := domain.NewPassword(passPlane)
	if err != nil {
		return nil, err
	}
	p, err := domain.NewPass(pID, aID, password, iat)
	if err != nil {
		return nil, err
	}

	// 保存
	_, err2 := a.userRepo.Insert(ctx, u)

	if err2 != nil {
		return nil, WrapUsecaseError("保存に失敗しました", err2)
	}
	_, err2 = a.passRepo.Insert(ctx, *p)
	if err2 != nil {
		return nil, WrapUsecaseError("保存に失敗しました", err2)
	}

	// プレゼン層に返す値を生成
	saveUserDto := NewUserDto(u.AccountID(), u.ID(), u.Name())
	return saveUserDto, nil
}

func (a *authN) SignIn(ctx context.Context, userID string, passwordPlane string) (bool, UsecaseError) {
	uID, err := domain.NewUserID(userID)
	if err != nil {
		return false, WrapUsecaseError("ログインに失敗しました", err)
	}

	u, err2 := a.userRepo.FindById(ctx, uID)
	if err2 != nil {
		return false, WrapUsecaseError("ログインに失敗しました", err2)
	}

	p, err2 := a.passRepo.FindByAccountId(ctx, u.AccountID())
	if err2 != nil {
		return false, WrapUsecaseError("ログインに失敗しました", err2)
	}

	// パスワードが一致するか
	isOK, err := p.IsLogin(passwordPlane)
	if err != nil {
		return false, WrapUsecaseError("ログインに失敗しました", err)
	}

	return isOK, nil
}

func (a *authN) SignOut(ctx context.Context, userID string) (bool, UsecaseError) {
	// TODO: ログアウトの処理とは?
	isOK := true
	return isOK, nil
}

func (a *authN) GetUserList(ctx context.Context) ([]UserDto, UsecaseError) {
	us, err := a.userRepo.Fetch(ctx)
	if err != nil {
		return nil, WrapUsecaseError("ユーザーの取得に失敗しました", err)
	}

	var userDtos []UserDto
	for _, u := range us {
		userDto := NewUserDto(u.AccountID(), u.ID(), u.Name())
		userDtos = append(userDtos, *userDto)
	}
	return userDtos, nil
}

func (a *authN) SearchUser(ctx context.Context, userID string) (*UserDto, UsecaseError) {
	uID, err := domain.NewUserID(userID)
	if err != nil {
		return nil, WrapUsecaseError("ユーザーの検索に失敗しました", err)
	}

	u, err2 := a.userRepo.FindById(ctx, uID)
	if err2 != nil {
		return nil, WrapUsecaseError("ユーザーの検索に失敗しました", err2)
	}

	userDto := NewUserDto(u.AccountID(), u.ID(), u.Name())
	return userDto, nil
}
