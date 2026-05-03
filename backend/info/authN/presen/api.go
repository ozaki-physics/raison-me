package presen

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ozaki-physics/raison-me/info/authN/usecase"
	// middleware は presen 層 と考えられるため 問題ない
	sharemiddleware "github.com/ozaki-physics/raison-me/share/middleware"
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

	// アクセストークン を 更新
	RefreshAccessToken(w http.ResponseWriter, r *http.Request) error
}

type apiCase struct {
	usecase.AuthN
	// TODO: ローカル開発で Cookie の Secure 属性 を 付けないようにするため
	// ローカルでも https で 開発するようにするのが 理想だけど とりあえず 制御 できるようにしている
	isLive bool
}

func NewAPICase(authn usecase.AuthN, isLive bool) ApiCase {
	return &apiCase{authn, isLive}
}

// TODO: 動作確認用で 暫定な実装
func (api *apiCase) CreateUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	type createUserRequest struct {
		UserName string `json:"userName"`
		UserID   string `json:"userID"`
		Password string `json:"password"`
	}
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return NewPresenError("invalid request body", http.StatusBadRequest)
	}
	if strings.TrimSpace(req.UserName) == "" || strings.TrimSpace(req.UserID) == "" || strings.TrimSpace(req.Password) == "" {
		return NewPresenError("userName, userID and password are required", http.StatusBadRequest)
	}

	saveUserDto, err := api.AuthN.Create(ctx, req.UserName, req.UserID, req.Password)
	if err != nil {
		return WrapPresenError("内部エラー", err, http.StatusInternalServerError)
	}

	return writeJSON(w, http.StatusCreated, saveUserDto)
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
		return NewPresenError("invalid request body", http.StatusBadRequest)
	}
	if strings.TrimSpace(req.UserID) == "" || strings.TrimSpace(req.Password) == "" {
		return NewPresenError("userID and password are required", http.StatusBadRequest)
	}

	result, err := api.AuthN.SignIn(ctx, req.UserID, req.Password)
	if err != nil {
		return WrapPresenError("内部エラー", err, http.StatusInternalServerError)
	}
	// リフレッシュトークン を Cookie にセットする
	api.setRefreshTokenCookie(w, result.RefreshToken)
	// アクセストークン を レスポンスボディ に返す JSON にする
	j := NewAccessTokenResponse(result.AccessToken)

	return writeJSON(w, http.StatusOK, j)
}

// TODO: 動作確認用で 暫定な実装
func (api *apiCase) SignOut(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	// Cookie から リフレッシュトークン を 取得する
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		return NewPresenError("refreshToken cookie is required", http.StatusBadRequest)
	}
	refreshToken := cookie.Value

	if err := api.AuthN.SignOut(ctx, refreshToken); err != nil {
		return WrapPresenError("内部エラー", err, http.StatusInternalServerError)
	}

	// サインアウト したあとは クライアント側の Cookie も 削除する
	api.clearRefreshTokenCookie(w)

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (api *apiCase) GetUserList(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	users, err := api.AuthN.GetUserList(ctx)
	if err != nil {
		return WrapPresenError("内部エラー", err, http.StatusInternalServerError)
	}

	return writeJSON(w, http.StatusOK, users)
}

// TODO: 動作確認用で 暫定な実装
func (api *apiCase) SearchUser(w http.ResponseWriter, r *http.Request, param string) error {
	ctx := r.Context()

	user, err := api.AuthN.SearchUser(ctx, param)
	if err != nil {
		return WrapPresenError("内部エラー", err, http.StatusInternalServerError)
	}

	return writeJSON(w, http.StatusOK, user)
}

// TODO: 動作確認用で 暫定な実装
func (api *apiCase) GetMe(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims, ok := sharemiddleware.AccessTokenClaimsFromContext(ctx)
	if !ok {
		return NewPresenError("failed to get access token claims from context", http.StatusBadRequest)
	}
	aID := claims.AccountID
	user, err := api.AuthN.GetMe(ctx, aID.Val())
	if err != nil {
		return WrapPresenError("内部エラー", err, http.StatusInternalServerError)
	}

	return writeJSON(w, http.StatusOK, user)
}

// TODO: 動作確認用で 暫定な実装
func (api *apiCase) RefreshAccessToken(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	// Cookie から リフレッシュトークン を 取得する
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		return NewPresenError("refreshToken cookie is required", http.StatusBadRequest)
	}
	refreshToken := cookie.Value

	result, err := api.AuthN.RefreshAccessToken(ctx, refreshToken)
	if err != nil {
		return WrapPresenError("内部エラー", err, http.StatusUnauthorized)
	}

	// 新しい リフレッシュトークン を Cookie にセットする
	api.setRefreshTokenCookie(w, result.RefreshToken)
	// 新しい アクセストークン を レスポンスボディ に返す JSON にする
	j := NewAccessTokenResponse(result.AccessToken)

	return writeJSON(w, http.StatusOK, j)
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// リフレッシュトークン を Cookie に セットするための ヘルパー関数
func (api *apiCase) setRefreshTokenCookie(w http.ResponseWriter, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		HttpOnly: true,
		// 本番環境 なら Secure 属性 を 付与する
		Secure: api.isLive,
		Path:   "/",
		// この設定にすると クロスサイト からの リクエスト には Cookie が 送信されなくなるため CSRF 対策になる
		SameSite: http.SameSiteStrictMode,
		// TODO: MaxAge の 設定は どうする?
		// infra 層に書かれている リフレッシュトークンの有効期限 と 同じがよいが
		// infra 層の実装に 依存 するのは よくない気がするので どうするか要検討
		MaxAge: 60 * 60 * 24 * 7, // 7日間
	})
}

// リフレッシュトークン を Cookie から 削除するための ヘルパー関数
// MaxAge を負の値にして 空の値をセットする
func (api *apiCase) clearRefreshTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		HttpOnly: true,
		// 本番環境 なら Secure 属性 を 付与する
		Secure: api.isLive,
		Path:   "/",
		// この設定にすると クロスサイト からの リクエスト には Cookie が 送信されなくなるため CSRF 対策になる
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
}
