package trysupabase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// postgREST API を 使った 接続確認
// apiKey は service_role キー を 使う
// Authorization ヘッダー を 付与
// 通常 は service_role キー は クライアント 側 で は 使わない こと
// HTTP クライアント を 作って 使う パターン
func getRows() {
	config, err := loadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	tableName := "tests"
	baseURL := config.SupabaseURL + "/rest/v1/" + tableName
	apiKey := config.SupabaseServiceRoleKey

	// SELECT (GET)
	req, _ := http.NewRequest("GET", baseURL+"?select=*", nil)
	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	fmt.Println("GET:", string(body))
}

// postgREST API を 使った 接続確認
// apiKey は secret キー を 使う
// Authorization ヘッダー は 付与しない
// HTTP クライアント は デフォルト を 使う パターン
func getRowsV2() {
	config, err := loadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	tableName := "tests"
	baseURL := config.SupabaseURL + "/rest/v1/" + tableName
	apiKey := config.SupabaseSecretKey
	// SELECT (GET)
	req, _ := http.NewRequest("GET", baseURL+"?select=*", nil)
	req.Header.Set("apikey", apiKey)
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	fmt.Println("GET v2:", string(body))
}

// postgREST API を 使った WHERE 付き 接続確認
// apiKey は service_role キー を 使う
// Authorization ヘッダー を 付与
// WHERE 句 は クエリパラメータ で 指定
// HTTP クライアント を 作って 使う パターン
func getColumnsWhere(id int) {
	config, err := loadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	tableName := "tests"
	baseURL := config.SupabaseURL + "/rest/v1/" + tableName
	apiKey := config.SupabaseServiceRoleKey

	// SELECT (GET)
	req, _ := http.NewRequest("GET", baseURL+"?select=*&id=eq."+fmt.Sprint(id), nil)
	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	fmt.Println("GET with WHERE:", string(body))
}

// INSERT 確認 用 の 構造体
type row struct {
	ID                   int    `json:"id"`
	TestText             string `json:"test_text"`
	TestBool             bool   `json:"test_bool"`
	TestTimestamptzNow   string `json:"test_timestamptz_now"`
	TestTimestampNowtz   string `json:"test_timestamp_nowtz"`
	TestTimestamptzNowtz string `json:"test_timestamptz_nowtz"`
	CreatedAt            string `json:"created_at"`
}

// postgREST API を 使った INSERT 確認
// apiKey は service_role キー を 使う
// Authorization ヘッダー を 付与
// HTTP クライアント を 作って 使う パターン
func insertRow(name string) {
	config, err := loadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	tableName := "tests"
	baseURL := config.SupabaseURL + "/rest/v1/" + tableName
	apiKey := config.SupabaseServiceRoleKey

	// INSERT (POST)
	newTest := row{
		ID:                   1,
		TestText:             name,
		TestBool:             false,
		TestTimestamptzNow:   "2025-10-10T15:21:07+00:00",
		TestTimestampNowtz:   "2025-10-11T00:21:08",
		TestTimestamptzNowtz: "2025-10-12T15:21:09+00:00",
		CreatedAt:            "2025-10-13T15:21:06+00:00",
	}
	data, _ := json.Marshal(newTest)

	req2, _ := http.NewRequest("POST", baseURL, bytes.NewReader(data))
	req2.Header.Set("apikey", apiKey)
	req2.Header.Set("Authorization", "Bearer "+apiKey)
	req2.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res2, err := client.Do(req2)
	if err != nil {
		panic(err)
	}
	defer res2.Body.Close()

	body2, _ := io.ReadAll(res2.Body)
	fmt.Println("POST:", string(body2))
}

// postgREST API を 使った スキーマ 取得 確認
// apiKey は secret キー を 使う
// Authorization ヘッダー は 付与しない
// ヘッダー に Accept-Profile を 付与
// HTTP クライアント は デフォルト を 使う パターン
func getRowsSchema() {
	config, err := loadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	tableName := "users"
	baseURL := config.SupabaseURL + "/rest/v1/" + tableName
	apiKey := config.SupabaseSecretKey
	// SELECT (GET)
	req, _ := http.NewRequest("GET", baseURL+"?select=*", nil)
	req.Header.Set("apikey", apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Profile", "app")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	fmt.Println("GET Schema:", string(body))
}

type appUsersRow struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// postgREST API を 使った INSERT 確認
// apiKey は service_role キー を 使う
// Authorization ヘッダー を 付与
// ヘッダー に Content-Profile を 付与
// HTTP クライアント を 作って 使う パターン
func insertAppUsersRow(name string) {
	config, err := loadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	tableName := "users"
	baseURL := config.SupabaseURL + "/rest/v1/" + tableName
	apiKey := config.SupabaseServiceRoleKey

	newRow := appUsersRow{
		Name: name,
	}

	// JSON.Marshal の 代わり に 手動で JSON 文字列 を 作成
	// なぜなら id が zero 値 で 入力されてしまい Supabase で 自動採番 されなくなるから
	jsonStr := `{"name":` + strconv.Quote(newRow.Name) + `}`
	data := []byte(jsonStr)

	// INSERT (POST)
	req, _ := http.NewRequest("POST", baseURL, bytes.NewReader(data))
	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Profile", "app")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	fmt.Println("INSERT App Users:", string(body))
}
