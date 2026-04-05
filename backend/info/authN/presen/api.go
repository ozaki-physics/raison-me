package presen

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ozaki-physics/raison-me/info/authN/usecase"
)

type ApiCase interface {
	// ユーザー で サインイン
	SignIn(w http.ResponseWriter, r *http.Request) error
	// ユーザー で サインアウト
	SignOut(w http.ResponseWriter, r *http.Request) error

	// ユーザー を 作成
	CreateUser(w http.ResponseWriter, r *http.Request) error
	// ユーザー一覧
	GetUserList(w http.ResponseWriter, r *http.Request) error
	// ユーザーを検索
	SearchUser(w http.ResponseWriter, r *http.Request, req string) error

	// 自分自身 の 情報を取得
	GetMe(w http.ResponseWriter, r *http.Request) error

	// IDトークン の 発行
	GenerateToken(w http.ResponseWriter, r *http.Request) error
	// IDトークン が 有効か
	VerifyToken(w http.ResponseWriter, r *http.Request) error
	// IDトークン の 更新
	RefreshToken(w http.ResponseWriter, r *http.Request) error
}

type apiCase struct {
	usecase.AuthN
}

func NewAPICase(authn usecase.AuthN) ApiCase {
	return &apiCase{authn}
}

// TODO: 動作確認用で 暫定な実装
func (api *apiCase) CreateUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	userName := r.FormValue("userName")
	userID := r.FormValue("userID")
	passPlane := r.FormValue("password")

	saveUserDto, err := api.AuthN.Create(ctx, userName, userID, passPlane)
	if err != nil {
		return err
	}

	w.Write([]byte("ユーザーを作成しました\n"))
	w.Write([]byte("AccountID: " + saveUserDto.AccountID + "\n"))
	w.Write([]byte("UserID: " + saveUserDto.UserID + "\n"))
	w.Write([]byte("UserName: " + saveUserDto.UserName + "\n"))
	return nil
}

// TODO: 動作確認用で 暫定な実装
func (api *apiCase) SignIn(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	type signInRequest struct {
		UserID   string `json:"userID"`
		Password string `json:"password"`
	}
	var req signInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	userID := req.UserID
	passwordPlane := req.Password

	isOK, err := api.AuthN.SignIn(ctx, userID, passwordPlane)
	if err != nil {
		return err
	}

	if isOK {
		w.Write([]byte("ログインに成功しました\n"))
	} else {
		w.Write([]byte("ログインに失敗しました\n"))
	}
	return nil
}

// TODO: 動作確認用で 暫定な実装
func (api *apiCase) SignOut(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("ログアウトに失敗しました\n"))
	return nil
}

// TODO: 動作確認用で 暫定な実装
func (api *apiCase) GetUserList(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	users, err := api.AuthN.GetUserList(ctx)
	if err != nil {
		return err
	}

	for _, u := range users {
		s := fmt.Sprintf("AccountID: %s, UserID: %s, UserName: %s", u.AccountID, u.UserID, u.UserName)
		w.Write([]byte(s + "\n"))
	}
	return nil
}

// TODO: 動作確認用で 暫定な実装
func (api *apiCase) SearchUser(w http.ResponseWriter, r *http.Request, param string) error {
	ctx := r.Context()

	user, err := api.AuthN.SearchUser(ctx, param)
	if err != nil {
		return err
	}

	s := fmt.Sprintf("AccountID: %s, UserID: %s, UserName: %s", user.AccountID, user.UserID, user.UserName)
	w.Write([]byte(s + "\n"))
	return nil
}

// TODO: 動作確認用で 暫定な実装
func (api *apiCase) GetMe(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("自分自身の情報を取得しました\n"))
	return nil
}

// TODO: 動作確認用で 暫定な実装
func (api *apiCase) GenerateToken(w http.ResponseWriter, r *http.Request) error {
	// TODO: 未実装
	fmt.Println("best")
	var jsonBody map[string]interface{}
	json.NewDecoder(r.Body).Decode(&jsonBody)

	fmt.Println(&jsonBody)

	return nil
}

// TODO: 動作確認用で 暫定な実装
func (api *apiCase) VerifyToken(w http.ResponseWriter, r *http.Request) error {
	// TODO: 未実装
	fmt.Println("lost")

	if err := r.ParseForm(); err != nil {
		return err
	}

	for k, v := range r.Form {
		fmt.Printf("%v: %v\n", k, v)
	}

	fmt.Println(r.FormValue("hello"))
	fmt.Println(r.PostFormValue("hello"))
	return nil
}

// TODO: 動作確認用で 暫定な実装
func (api *apiCase) RefreshToken(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("refresh")
	return nil
}
