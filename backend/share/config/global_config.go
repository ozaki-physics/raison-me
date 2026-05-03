package config

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Config interface {
	IsLive() bool
	IsCloud() bool
	GetGCPProjectID() string
	GetDSN() string
	GetSupabaseDSN() string
	GetPort() string
	// TODO: 暫定の認証
	GetSampleAPIToken() string
	GetAuthNPepper() string
	GetAuthNJWTSecret() string
	GetAuthNRefreshTokenPepper() string
}

func NewConfig() Config {
	log.Println("Config: called")
	isCloud := os.Getenv("IS_CLOUD") == "true"
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	isLive := readFile[string](isCloud, "IS_LIVE") == "true"
	gcpProjectID := readFile[string](isCloud, "GCP_PROJECT_ID")

	dataSourceName := readFile[string](isCloud, "DATABASE_URL")
	supabaseConfig := readFile[supabaseConfig](isCloud, "DATABASE_SUPABASE_JSON")

	// TODO: 暫定の認証
	sampleAPIToken := readFile[string](isCloud, "SAMPLE_API_TOKEN")

	authNPepper := readFile[string](isCloud, "AUTHN_PASSWORD_PEPPER")
	authNJWTSecret := readFile[string](isCloud, "AUTHN_JWT_HS256_SECRET")
	authNRefreshTokenPepper := readFile[string](isCloud, "AUTHN_REFRESH_TOKEN_PEPPER")

	c := config{
		isLive:                  isLive,
		isCloud:                 isCloud,
		gcpProjectID:            gcpProjectID,
		dataSourceName:          dataSourceName,
		supabaseConfig:          supabaseConfig,
		port:                    port,
		sampleAPIToken:          sampleAPIToken,
		authNPepper:             authNPepper,
		authNJWTSecret:          authNJWTSecret,
		authNRefreshTokenPepper: authNRefreshTokenPepper,
	}
	return &c
}

type config struct {
	isLive                  bool
	isCloud                 bool
	gcpProjectID            string
	dataSourceName          string
	supabaseConfig          supabaseConfig
	port                    string
	sampleAPIToken          string
	authNPepper             string
	authNJWTSecret          string
	authNRefreshTokenPepper string
}

func (c *config) IsLive() bool {
	return c.isLive
}

func (c *config) IsCloud() bool {
	return c.isCloud
}

func (c *config) GetGCPProjectID() string {
	return c.gcpProjectID
}

func (c *config) GetDSN() string {
	return c.dataSourceName
}

func (c *config) GetSupabaseDSN() string {
	config := c.supabaseConfig
	escapedPass := url.QueryEscape(config.Password)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=require", config.User, escapedPass, config.Host, config.TransactionPort, config.Dbname)
	return dsn
}

func (c *config) GetPort() string {
	return c.port
}

func (c *config) GetSampleAPIToken() string {
	return c.sampleAPIToken
}

func (c *config) GetAuthNPepper() string {
	return c.authNPepper
}

func (c *config) GetAuthNJWTSecret() string {
	return c.authNJWTSecret
}

func (c *config) GetAuthNRefreshTokenPepper() string {
	return c.authNRefreshTokenPepper
}

type supabaseConfig struct {
	User            string `json:"user"`
	Password        string `json:"password"`
	Host            string `json:"host"`
	TransactionPort int    `json:"transaction_port"`
	Dbname          string `json:"dbname"`
}

// 機密情報 JSON ファイル を 読み取る ヘルパー 関数
// (ジェネリクス 使用のため レシーバー にできない)
func readFile[T any](isCloud bool, fileName string) T {
	// ファイル パス の 組み立て
	var filePath string
	if isCloud {
		// なぜ fileName を 2回 繰り返すか
		// Cloud Run では Secret Manager は 1つの シークレット に 1つの ボリューム を 対応 させる 必要がある
		// また ボリューム を 1個のディレクトリに まとめられない
		// よって シークレット名 ごとに ディレクトリ が 作成 されるため

		// /app は Cloud Run 用 の Dockerfile で 作成 した ディレクトリ
		filePath = filepath.Join("/app/share/secrets/", fileName, fileName)
	} else {
		// secret file が 格納されてる ディレクトリ の パスを取得
		filePath = filepath.Join("./share/secrets/", fileName)
	}

	var zero T
	if _, ok := any(zero).(string); ok {
		b, err := os.ReadFile(filePath)
		if err != nil {
			// TODO: エラーハンドリング が 雑
			log.Printf("Failed to read file: %v", err)
			return zero
		}

		s := strings.TrimSpace(string(b))
		return any(s).(T)
	}

	b, err := os.ReadFile(filePath)
	if err != nil {
		// TODO: エラーハンドリング が 雑
		log.Fatalf("Failed to read file(JSON): %v", err)
	}

	var result T
	err = json.Unmarshal(b, &result)
	if err != nil {
		// TODO: エラーハンドリング が 雑
		log.Fatalf("Failed to parse JSON file: %v", err)
	}

	return result
}
