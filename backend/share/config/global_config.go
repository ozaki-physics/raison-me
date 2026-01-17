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
}

func NewConfig() Config {
	isCloud := os.Getenv("IS_CLOUD") == "true"
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	sc := newSecretConfig(isCloud)
	isLive := sc.readFile("IS_LIVE") == "true"
	gcpProjectID := sc.readFile("GCP_PROJECT_ID")

	dataSourceName := sc.readFile("DATABASE_URL")
	supabaseConfig := readJSONFile[supabaseConfig](sc, "DATABASE_SUPABASE_JSON")

	c := config{
		isLive:         isLive,
		isCloud:        isCloud,
		gcpProjectID:   gcpProjectID,
		dataSourceName: dataSourceName,
		supabaseConfig: supabaseConfig,
		port:           port,
	}
	return &c
}

type config struct {
	isLive         bool
	isCloud        bool
	gcpProjectID   string
	dataSourceName string
	supabaseConfig supabaseConfig
	port           string
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

// secretConfig は secret file の 読み取り に関する 設定 を 管理
type secretConfig struct {
	basePath string
}

func newSecretConfig(isCloud bool) *secretConfig {
	// secret file が 格納されてる ディレクトリ の パスを取得
	basePath := "./share/secrets/"
	if isCloud {
		// /app は Cloud Run 用 の Dockerfile で 作成 した ディレクトリ
		basePath = "/app/share/secrets/"
	}

	return &secretConfig{
		basePath: basePath,
	}
}

// 秘密情報 ファイル を 読み取る ヘルパー 関数
func (sc *secretConfig) readFile(fileName string) string {
	filePath := sc.basePath + fileName
	b, err := os.ReadFile(filePath)
	if err != nil {
		// TODO: エラーハンドリング が 雑
		log.Printf("Failed to read file: %v", err)
		return ""
	}

	value := string(b)
	// 改行コード を 削除
	value = strings.TrimSpace(value)
	return value
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
func readJSONFile[T any](sc *secretConfig, fileName string) T {
	filePath := filepath.Join(sc.basePath, fileName)
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
