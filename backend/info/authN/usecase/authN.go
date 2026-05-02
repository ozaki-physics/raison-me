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
	SignIn(ctx context.Context, userID string, passwordPlane string) (*AuthResultDto, UsecaseError)
	// ユーザー で サインアウト
	SignOut(ctx context.Context, refreshToken string) UsecaseError
	// リフレッシュトークン で アクセストークン を 更新
	RefreshAccessToken(ctx context.Context, refreshToken string) (*AuthResultDto, UsecaseError)
	// ユーザー(自分) の情報を更新
	// Update(userDto UserDto) (*UserDto, UsecaseError)
	// パスワードリセット
	// ResetPassword()
	// ユーザー一覧 を 取得
	GetUserList(ctx context.Context) ([]UserDto, UsecaseError)
	// ユーザー を 検索
	SearchUser(ctx context.Context, userID string) (*UserDto, UsecaseError)
	// 自分 の 情報 を 取得
	GetMe(ctx context.Context, accountID string) (*UserDto, UsecaseError)
}

type authN struct {
	userRepo              domain.UserRepo
	passRepo              domain.PassRepo
	refreshTokenRepo      domain.RefreshTokenRepo
	hasher                domain.PasswordHasherRepo
	accessTokenIssuer     AccessTokenIssuer
	refreshTokenGenerator RefreshTokenIssuer
	refreshTokenHasher    RefreshTokenHasher
}

func NewAuthN(
	u domain.UserRepo,
	p domain.PassRepo,
	rt domain.RefreshTokenRepo,
	h domain.PasswordHasherRepo,
	ati AccessTokenIssuer,
	rtg RefreshTokenIssuer,
	rth RefreshTokenHasher,
) (AuthN, UsecaseError) {
	if u == nil {
		return nil, NewUsecaseError("UserRepo が設定されていません")
	}
	if p == nil {
		return nil, NewUsecaseError("PassRepo が設定されていません")
	}
	if rt == nil {
		return nil, NewUsecaseError("RefreshTokenRepo が設定されていません")
	}
	if h == nil {
		return nil, NewUsecaseError("PasswordHasher が設定されていません")
	}
	if ati == nil {
		return nil, NewUsecaseError("AccessTokenIssuer が設定されていません")
	}
	if rtg == nil {
		return nil, NewUsecaseError("RefreshTokenGenerator が設定されていません")
	}
	if rth == nil {
		return nil, NewUsecaseError("RefreshTokenHasher が設定されていません")
	}

	a := &authN{
		userRepo:              u,
		passRepo:              p,
		refreshTokenRepo:      rt,
		hasher:                h,
		accessTokenIssuer:     ati,
		refreshTokenGenerator: rtg,
		refreshTokenHasher:    rth,
	}
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
	passHash, err3 := a.hasher.Hash(passPlane)
	if err3 != nil {
		return nil, WrapUsecaseError("パスワードの生成に失敗しました", err3)
	}
	password, err := domain.NewPassword(passHash)
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

// SignIn は ユーザーID と パスワード を 受け取って 認証する
func (a *authN) SignIn(ctx context.Context, userID string, passwordPlane string) (*AuthResultDto, UsecaseError) {
	uID, derr := domain.NewUserID(userID)
	if derr != nil {
		return nil, WrapUsecaseError("ログインに失敗しました", derr)
	}

	u, err := a.userRepo.FindById(ctx, uID)
	if err != nil {
		return nil, WrapUsecaseError("ログインに失敗しました", err)
	}

	p, err := a.passRepo.FindByAccountId(ctx, u.AccountID())
	if err != nil {
		return nil, WrapUsecaseError("ログインに失敗しました", err)
	}
	// パスワードが一致するか
	if err := a.hasher.Verify(p.Password(), passwordPlane); err != nil {
		return nil, WrapUsecaseError("ログインに失敗しました", err)
	}

	accessToken, err := a.accessTokenIssuer.Issue(u.AccountID())
	if err != nil {
		return nil, WrapUsecaseError("ログインに失敗しました", err)
	}
	refreshToken, err := a.refreshTokenGenerator.Issue()
	if err != nil {
		return nil, WrapUsecaseError("ログインに失敗しました", err)
	}
	refreshTokenHash, err := a.refreshTokenHasher.Hash(refreshToken.Token())
	if err != nil {
		return nil, WrapUsecaseError("ログインに失敗しました", err)
	}

	rID, err := domain.NewRefreshTokenID()
	if err != nil {
		return nil, WrapUsecaseError("ログインに失敗しました", err)
	}
	now, err := domain.NewDate()
	if err != nil {
		return nil, WrapUsecaseError("ログインに失敗しました", err)
	}
	refreshTokenClaims, err := domain.NewRefreshTokenClaims(rID, u.AccountID(), refreshTokenHash, refreshToken.ExpiresAt(), now)
	if err != nil {
		return nil, WrapUsecaseError("ログインに失敗しました", err)
	}

	_, err = a.refreshTokenRepo.Insert(ctx, refreshTokenClaims)
	if err != nil {
		return nil, WrapUsecaseError("ログインに失敗しました", err)
	}

	authResultDto := NewAuthResultDto(accessToken.Token(), refreshToken.Token())
	return authResultDto, nil
}

// SignOut は リフレッシュトークン を 受け取って ログアウトする
func (a *authN) SignOut(ctx context.Context, refreshToken string) UsecaseError {
	// TODO: リフレッシュトークン の 扱いだけでいいの? ログアウトに JWT の 扱いは?
	hashedRefreshToken, err := a.refreshTokenHasher.Hash(refreshToken)
	if err != nil {
		return WrapUsecaseError("ログアウトに失敗しました", err)
	}

	refreshTokenClaims, err := a.refreshTokenRepo.FindByTokenHash(ctx, hashedRefreshToken)
	if err != nil {
		return WrapUsecaseError("ログアウトに失敗しました", err)
	}

	// リフレッシュトークンが既に失効している場合は 何もしない
	if refreshTokenClaims.IsRevoked() {
		return nil
	}

	// リフレッシュトークン を 失効させる
	err4 := a.refreshTokenRepo.Revoke(ctx, hashedRefreshToken)
	if err4 != nil {
		return WrapUsecaseError("ログアウトに失敗しました", err4)
	}

	return nil
}

// RefreshAccessToken は リフレッシュトークン を 受け取って アクセストークン と リフレッシュトークン を 更新する
func (a *authN) RefreshAccessToken(ctx context.Context, refreshToken string) (*AuthResultDto, UsecaseError) {
	hashedRefreshToken, err := a.refreshTokenHasher.Hash(refreshToken)
	if err != nil {
		return nil, WrapUsecaseError("リフレッシュトークンの処理に失敗しました", err)
	}
	// 現時点の ハッシュ化された リフレッシュトークン をもとに 現時点の クレーム を 取得する
	refreshTokenClaims, err := a.refreshTokenRepo.FindByTokenHash(ctx, hashedRefreshToken)
	if err != nil {
		return nil, WrapUsecaseError("リフレッシュトークンの処理に失敗しました", err)
	}

	// 現時点の リフレッシュトークン が 有効か確認
	d, derr := domain.NewDate()
	if derr != nil {
		return nil, WrapUsecaseError("リフレッシュトークンの処理, 日時生成に失敗しました", derr)
	}
	if !refreshTokenClaims.IsActive(d) {
		return nil, NewUsecaseError("リフレッシュトークンは無効です")
	}

	// 新しい リフレッシュトークン を 作り ハッシュ化 する
	nextRefreshToken, err := a.refreshTokenGenerator.Issue()
	if err != nil {
		return nil, WrapUsecaseError("リフレッシュトークンの処理に失敗しました", err)
	}
	nextRefreshTokenHash, err := a.refreshTokenHasher.Hash(nextRefreshToken.Token())
	if err != nil {
		return nil, WrapUsecaseError("リフレッシュトークンの処理に失敗しました", err)
	}

	// クレーム を 更新して 新しい クレーム を 作る
	rotatedRefreshToken, err := refreshTokenClaims.WithRotatedToken(nextRefreshTokenHash, nextRefreshToken.ExpiresAt(), d)
	if err != nil {
		return nil, WrapUsecaseError("リフレッシュトークンの処理に失敗しました", err)
	}

	// データストア上の リフレッシュトークン を 更新する
	updatedRefreshToken, err := a.refreshTokenRepo.Rotate(ctx, hashedRefreshToken, rotatedRefreshToken)
	if err != nil {
		return nil, WrapUsecaseError("リフレッシュトークンの処理に失敗しました", err)
	}

	// 新しい アクセストークン を 発行する
	accessToken, err := a.accessTokenIssuer.Issue(updatedRefreshToken.AccountID())
	if err != nil {
		return nil, WrapUsecaseError("アクセストークンの発行に失敗しました", err)
	}

	return NewAuthResultDto(accessToken.Token(), nextRefreshToken.Token()), nil
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

// GetMe は 認証されたユーザー 自身 の 情報 を 取得する
func (a *authN) GetMe(ctx context.Context, accountID string) (*UserDto, UsecaseError) {
	aID, derr := domain.ReNewAccountID(accountID)
	if derr != nil {
		return nil, WrapUsecaseError("自身の情報の取得に失敗しました", derr)
	}
	u, err := a.userRepo.FindByAccountId(ctx, aID)
	if err != nil {
		return nil, WrapUsecaseError("自身の情報の取得に失敗しました", err)
	}

	return NewUserDto(u.AccountID(), u.ID(), u.Name()), nil
}
