package presentation

import (
	"fmt"
	"net/http"
)

// Routes プレゼンテーションのルータのインタフェース
type Routes interface {
	Handler(w http.ResponseWriter, r *http.Request)
}

// routes プレゼンテーションのルータの実体
type routes struct {
	caPresen CryptoAssetsPresen
}

// CreateCryptoAssetsRoutes プレゼンテーションのルータのコンストラクタ
// 戻り値がインタフェースだから 実装を強制できる
func CreateCryptoAssetsRoutes(caPresen CryptoAssetsPresen) Routes {
	return &routes{caPresen}
}

// Handler routes の初期化
// http.HandleFunc の引数に渡せるのは http.ResponseWriter, *http.Request だけを引数に持ったメソッド
func (ro *routes) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		fmt.Fprintln(w, "\nGET だよ")
		err := ro.caPresen.Get(w, r)
		if err != nil {
			http.Error(w, "Internal Server Error\n"+err.Error(), 500)
		}
	// case http.MethodPost:
	// 	fmt.Fprintln(w, "\nPOST だよ")
	// 	err := ro.caPresen.Post(w, r)
	// 	if err != nil {
	// 		http.Error(w, "Internal Server Error\n"+err.Error(), 500)
	// 	}
	// case http.MethodPut:
	// 	fmt.Fprintln(w, "\nPUT だよ")
	// 	err := ro.caPresen.Put(w, r)
	// 	if err != nil {
	// 		http.Error(w, "Internal Server Error\n"+err.Error(), 500)
	// 	}
	// case http.MethodDelete:
	// 	fmt.Fprintln(w, "\nDELETE だよ")
	// 	err := ro.caPresen.Delete(w, r)
	// 	if err != nil {
	// 		http.Error(w, "Internal Server Error\n"+err.Error(), 500)
	// 	}
	default:
		fmt.Fprintln(w, "許可されたメソッドじゃないよ")
	}
}
