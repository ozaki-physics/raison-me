package trysupabase

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type supabaseConfig struct {
	SupabaseURL            string `json:"supabase_url"`
	SupabaseServiceRoleKey string `json:"supabase_service_role_key"`
	SupabaseSecretKey      string `json:"supabase_secret_key"`
	User                   string `json:"user"`
	Password               string `json:"password"`
	Host                   string `json:"host"`
	SessionPort            int    `json:"session_port"`
	TransactionPort        int    `json:"transaction_port"`
	DBName                 string `json:"dbname"`
}

func loadConfig() (*supabaseConfig, error) {
	data, err := os.ReadFile(filepath.Join("trysupabase", "key.json"))
	if err != nil {
		return nil, err
	}

	var config supabaseConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
