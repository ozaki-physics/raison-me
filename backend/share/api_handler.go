package share

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ozaki-physics/raison-me/share/errorer"
)

type Handler interface {
	// ServeHTTP(w http.ResponseWriter, r *http.Request)
	Handler(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	ac apiCase
}

type apiCase func(w http.ResponseWriter, r *http.Request) error

func NewApiHandler(ac apiCase) Handler {
	return &handler{ac}
}

// レスポンス を 常に JSON にするため
// TODO: 大局的 エラー キャッチのため?
func (h *handler) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := h.ac(w, r)
	if err != nil {
		newErrorResponse(w, err)
	}
}

type errorResponse struct {
	Status int    `json:"statusCode"`
	Msg    string `json:"error"`
}

// StatusCoder を 実装している エラー なら そこから ステータスコード を 取得するため
type StatusCoder interface {
	StatusCode() int
}

// 内部 err を クライアント に 返す形に整える
func newErrorResponse(w http.ResponseWriter, e error) http.ResponseWriter {
	var s int
	// e が StatusCoder を 実装しているなら そこから ステータスコード を 取得する
	if statusErr, ok := e.(StatusCoder); ok {
		s = statusErr.StatusCode()
	} else {
		s = http.StatusInternalServerError
	}

	// e が FullErrorer を 実装しているなら FullError を ログに出力する
	if fullErr, ok := e.(errorer.FullErrorer); ok {
		log.Printf("エラー発生: %s", fullErr.FullError())
	}

	w.WriteHeader(s)
	_ = json.NewEncoder(w).Encode(errorResponse{Status: s, Msg: e.Error()})
	return w
}
