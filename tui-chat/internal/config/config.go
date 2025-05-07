package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config はアプリケーション設定
type Config struct {
	Username     string `json:"username"`
	DefaultPort  string `json:"default_port"`
	DefaultHost  string `json:"default_host"`
	ColorEnabled bool   `json:"color_enabled"`
}

// DefaultConfig はデフォルト設定
var DefaultConfig = Config{
	Username:     "ユーザー",
	DefaultPort:  "8080",
	DefaultHost:  "localhost",
	ColorEnabled: true,
}

// LoadConfig は設定をファイルから読み込むメソッド
func LoadConfig(path string) (Config, error) {
	// デフォルト設定
	config := DefaultConfig

	// ファイルが存在しない場合はデフォルト設定を保存して返す
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// ディレクトリが存在しない場合は作成
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return config, err
		}

		// デフォルト設定を保存
		if err := SaveConfig(path, config); err != nil {
			return config, err
		}

		return config, nil
	}

	// ファイルを開く
	file, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer file.Close()

	// JSONデコード
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}

// SaveConfig は設定をファイルに保存するメソッド
func SaveConfig(path string, config Config) error {
	// ディレクトリが存在しない場合は作成
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// ファイルを作成
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// JSONエンコード（整形あり）
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return err
	}

	return nil
}

// GetConfigPath は設定ファイルのパスを取得するメソッド
func GetConfigPath() (string, error) {
	// ホームディレクトリを取得
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// 設定ディレクトリのパス
	configDir := filepath.Join(homeDir, ".config", "tui-chat")

	// 設定ファイルのパス
	configPath := filepath.Join(configDir, "config.json")

	return configPath, nil
}
