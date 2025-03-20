package presen

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ozaki-physics/raison-me/info/authN/usecase"
)

type ApiCase interface {
	// ユーザー を 作成
	UserCreate(w http.ResponseWriter, r *http.Request) error
	// ユーザー で ログイン
	UserEntry(w http.ResponseWriter, r *http.Request) error
	// IDトークン の 発行
	IDTokenGenerate(w http.ResponseWriter, r *http.Request) error
	// IDトークン が 有効か
	IsIDTokenOK(w http.ResponseWriter, r *http.Request) error
}

type apiCase struct {
	user       usecase.User
	credential usecase.Credential
}

func NewAPICase(u usecase.User, c usecase.Credential) ApiCase {
	return &apiCase{u, c}
}

func (api *apiCase) UserCreate(w http.ResponseWriter, r *http.Request) error {
	// TODO: 未実装
	fmt.Println("hello")
	return nil
}

func (api *apiCase) UserEntry(w http.ResponseWriter, r *http.Request) error {
	// TODO: 未実装
	fmt.Println("world")

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

func (api *apiCase) IDTokenGenerate(w http.ResponseWriter, r *http.Request) error {
	// TODO: 未実装
	fmt.Println("best")
	var jsonBody map[string]interface{}
	json.NewDecoder(r.Body).Decode(&jsonBody)

	fmt.Println(&jsonBody)

	return nil
}

func (api *apiCase) IsIDTokenOK(w http.ResponseWriter, r *http.Request) error {
	// TODO: 未実装
	fmt.Println("lost")
	return nil
}
