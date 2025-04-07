package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Generator GeneratorConfig `yaml:"generator"`
	MCP       MCPConfig       `yaml:"mcp"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type GeneratorConfig struct {
	DefaultLength int  `yaml:"defaultLength"`
	MinLength     int  `yaml:"minLength"`
	MaxLength     int  `yaml:"maxLength"`
	UseUppercase  bool `yaml:"useUppercase"`
	UseLowercase  bool `yaml:"useLowercase"`
	UseDigits     bool `yaml:"useDigits"`
	UseSpecial    bool `yaml:"useSpecial"`
}

type MCPConfig struct {
	RequireConsent bool `yaml:"requireConsent"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

func validate(cfg *Config) error {
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", cfg.Server.Port)
	}

	if cfg.Generator.MinLength < 1 {
		return fmt.Errorf("minimum password length must be at least 1")
	}

	if cfg.Generator.MaxLength < cfg.Generator.MinLength {
		return fmt.Errorf("maximum password length must be greater than or equal to minimum length")
	}

	if cfg.Generator.DefaultLength < cfg.Generator.MinLength || cfg.Generator.DefaultLength > cfg.Generator.MaxLength {
		return fmt.Errorf("default password length must be between min and max length")
	}

	if !cfg.Generator.UseUppercase && !cfg.Generator.UseLowercase &&
		!cfg.Generator.UseDigits && !cfg.Generator.UseSpecial {
		return fmt.Errorf("at least one character type must be enabled")
	}

	return nil
}
