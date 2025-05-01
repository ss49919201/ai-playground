package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// Config はアプリケーション全体の設定を保持します。
type Config struct {
	Server ServerConfig `yaml:"server"`
	Token  TokenConfig  `yaml:"token"`
	// Storage StorageConfig `yaml:"storage"` // 将来の拡張用
	// Crypto CryptoConfig `yaml:"crypto"` // 将来の拡張用
}

// ServerConfig はHTTPサーバー関連の設定を保持します。
type ServerConfig struct {
	Port        int    `yaml:"port"`
	TLSCertFile string `yaml:"tlsCertFile"` // TLS証明書ファイルパス
	TLSKeyFile  string `yaml:"tlsKeyFile"`  // TLS秘密鍵ファイルパス
}

// TokenConfig はトークン関連の設定（有効期間など）を保持します。
type TokenConfig struct {
	AccessTokenLifetime  time.Duration `yaml:"accessTokenLifetime"`
	RefreshTokenLifetime time.Duration `yaml:"refreshTokenLifetime"`
	AuthCodeLifetime     time.Duration `yaml:"authCodeLifetime"`
	// JWTSigningKeyFile string        `yaml:"jwtSigningKeyFile"` // JWT署名鍵ファイルパス
	// JWTIssuer         string        `yaml:"jwtIssuer"`         // JWT発行者
}

/*
// StorageConfig はストレージ関連の設定を保持します。
type StorageConfig struct {
	Type     string         `yaml:"type"` // "memory", "file", "database" など
	File     FileStorageConfig `yaml:"file"`
	Database DBStorageConfig   `yaml:"database"`
}

type FileStorageConfig struct {
	Path string `yaml:"path"`
}

type DBStorageConfig struct {
	DSN string `yaml:"dsn"` // Data Source Name
}

// CryptoConfig は暗号化関連の設定を保持します。
type CryptoConfig struct {
	// 例: パスワードハッシュ化のコストなど
	PasswordHashCost int `yaml:"passwordHashCost"`
}
*/

// Load は指定されたパスから設定ファイルを読み込み、解析して Config オブジェクトを返します。
// ファイルが存在しない、または解析に失敗した場合はエラーを返します。
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("設定ファイル '%s' の読み込みに失敗しました: %w", path, err)
	}

	// デフォルト値の設定
	cfg := Config{
		Server: ServerConfig{Port: 8080}, // デフォルトポート
		Token: TokenConfig{
			AccessTokenLifetime:  time.Hour * 1,       // デフォルト1時間
			RefreshTokenLifetime: time.Hour * 24 * 30, // デフォルト30日
			AuthCodeLifetime:     time.Minute * 10,    // デフォルト10分
		},
		/*
			Storage: StorageConfig{
				Type: "memory", // デフォルトはインメモリ
			},
			Crypto: CryptoConfig{
				PasswordHashCost: 0, // bcryptのデフォルトコストを使用
			},
		*/
	}

	// YAMLファイルの内容でデフォルト値を上書き
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("設定ファイル '%s' の解析に失敗しました: %w", path, err)
	}

	// 設定値のバリデーション
	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("設定が無効です: %w", err)
	}

	return &cfg, nil
}

// validate は Config オブジェクトの内容を検証します。
// 無効な値が含まれている場合はエラーを返します。
func validate(cfg *Config) error {
	// Server設定の検証
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("サーバーポート番号が無効です: %d", cfg.Server.Port)
	}
	if (cfg.Server.TLSCertFile != "" && cfg.Server.TLSKeyFile == "") || (cfg.Server.TLSCertFile == "" && cfg.Server.TLSKeyFile != "") {
		return fmt.Errorf("TLSを使用する場合、証明書ファイルと秘密鍵ファイルの両方を指定する必要があります")
	}

	// Token設定の検証
	if cfg.Token.AccessTokenLifetime <= 0 {
		return fmt.Errorf("アクセストークンの有効期間は正の値である必要があります: %v", cfg.Token.AccessTokenLifetime)
	}
	if cfg.Token.RefreshTokenLifetime <= 0 {
		return fmt.Errorf("リフレッシュトークンの有効期間は正の値である必要があります: %v", cfg.Token.RefreshTokenLifetime)
	}
	if cfg.Token.AuthCodeLifetime <= 0 {
		return fmt.Errorf("認可コードの有効期間は正の値である必要があります: %v", cfg.Token.AuthCodeLifetime)
	}
	// 一般的にリフレッシュトークンはアクセストークンより長い
	if cfg.Token.AccessTokenLifetime >= cfg.Token.RefreshTokenLifetime {
		// 警告を出すか、エラーにするかはポリシーによる
		// log.Printf("警告: アクセストークンの有効期間 (%v) がリフレッシュトークンの有効期間 (%v) 以上です", cfg.Token.AccessTokenLifetime, cfg.Token.RefreshTokenLifetime)
	}

	/*
		// Storage設定の検証 (将来の拡張用)
		switch cfg.Storage.Type {
		case "memory":
			// OK
		case "file":
			if cfg.Storage.File.Path == "" {
				return fmt.Errorf("ファイルストレージを使用する場合、パスを指定する必要があります")
			}
		case "database":
			if cfg.Storage.Database.DSN == "" {
				return fmt.Errorf("データベースストレージを使用する場合、DSNを指定する必要があります")
			}
		default:
			return fmt.Errorf("不明なストレージタイプです: %s", cfg.Storage.Type)
		}

		// Crypto設定の検証 (将来の拡張用)
		if cfg.Crypto.PasswordHashCost < 0 {
			return fmt.Errorf("パスワードハッシュコストは負の値にできません: %d", cfg.Crypto.PasswordHashCost)
		}
	*/

	return nil
}
