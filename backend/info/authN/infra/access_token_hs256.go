package infra

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/ozaki-physics/raison-me/info/authN/domain"
	"github.com/ozaki-physics/raison-me/info/authN/usecase"
)

// アクセストークン の 有効期限 (15分)
const accessTokenTTLSeconds = 900

type accessTokenHMAC struct {
	secret []byte
}

type jwtHeaderClaims struct {
	// 署名アルゴリズム を 表す (例: "HS256")
	Alg string `json:"alg"`
	// トークンの種類 を 表す (例: "JWT")
	Typ string `json:"typ"`
}

type jwtPayloadClaims struct {
	// トークンの主体 を 表す
	Sub string `json:"sub"`
	// Unix タイムスタンプ (秒) で 発行時間 を 表す
	Iat int64 `json:"iat"`
	// Unix タイムスタンプ (秒) で 有効期限 を 表す
	Exp int64 `json:"exp"`
}

// HMAC を 使った アクセストークン の 発行 と 検証
// HMAC とは 共通の秘密鍵 を 使って データ の 完全性 と 認証 を 確保する 暗号学的な手法
func NewAccessTokenHMAC(secret string) (*accessTokenHMAC, error) {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return nil, NewInfraError("署名する秘密鍵 は 必須です")
	}
	s := []byte(secret)

	hmac := &accessTokenHMAC{
		secret: s,
	}
	return hmac, nil
}

// アクセストークン を 発行
func (a *accessTokenHMAC) Issue(accountID domain.AccountID) (*usecase.AccessToken, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(time.Duration(accessTokenTTLSeconds) * time.Second)

	jh := jwtHeaderClaims{
		Alg: "HS256",
		Typ: "JWT",
	}
	jp := jwtPayloadClaims{
		Sub: accountID.Val(),
		Iat: now.Unix(),
		Exp: expiresAt.Unix(),
	}

	headerJSON, err := json.Marshal(jh)
	if err != nil {
		return nil, err
	}
	payloadJSON, err := json.Marshal(jp)
	if err != nil {
		return nil, err
	}

	header := base64.RawURLEncoding.EncodeToString(headerJSON)
	payload := base64.RawURLEncoding.EncodeToString(payloadJSON)
	signingInput := header + "." + payload
	signature := a.sign(signingInput)

	atr, err2 := usecase.ReNewAccessToken(signingInput+"."+signature, expiresAt)
	if err2 != nil {
		return nil, err2
	}
	return atr, nil
}

// アクセストークン を 検証
func (a *accessTokenHMAC) Verify(token string) (*domain.AccessTokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, NewInfraError("invalid token format")
	}
	header := parts[0]
	payload := parts[1]
	signature := parts[2]

	// TODO: トークンの各項目が 空じゃないことを確認しなくてよい?

	// 署名 を 検証する
	signingInput := header + "." + payload
	expectedSignature := a.sign(signingInput)
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		log.Printf("JWT の署名が違う: %s", token)
		return nil, NewInfraError("JWT の署名が違います")
	}

	// ペイロード を デコードして クレーム を 取得する
	payloadJSON, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return nil, WrapInfraError("failed to decode token payload: %w", err)
	}

	// ペイロード の クレーム を パースして アクセストークンのクレーム に 変換する
	var payloadClaims jwtPayloadClaims
	if err := json.Unmarshal(payloadJSON, &payloadClaims); err != nil {
		return nil, WrapInfraError("failed to parse token payload: %w", err)
	}
	if strings.TrimSpace(payloadClaims.Sub) == "" {
		return nil, NewInfraError("トークン の ペイロード が 空です")
	}

	// トークンの有効期限 を チェックする
	now := time.Now().UTC()
	expiresAt := time.Unix(payloadClaims.Exp, 0).UTC()
	if now.After(expiresAt) {
		log.Println("JWT の有効期限が切れている")
		return nil, NewInfraError("JWT の有効期限が切れています")
	}

	// アクセストークンのクレーム を 生成して返す
	atc, err02 := domain.ReNewAccessTokenClaims(
		payloadClaims.Sub,
		time.Unix(payloadClaims.Iat, 0).UTC(),
		expiresAt,
	)
	if err02 != nil {
		return nil, err02
	}

	return atc, nil
}

// HMAC で トークン を 署名する
func (a *accessTokenHMAC) sign(input string) string {
	mac := hmac.New(sha256.New, a.secret)
	mac.Write([]byte(input))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return signature
}
