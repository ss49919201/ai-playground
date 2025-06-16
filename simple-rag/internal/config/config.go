package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents the application configuration
type Config struct {
	Server struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`

	Embedding struct {
		URL       string `yaml:"url"`
		Model     string `yaml:"model"`
		BatchSize int    `yaml:"batch_size"`
	} `yaml:"embedding"`

	LLM struct {
		URL         string  `yaml:"url"`
		Model       string  `yaml:"model"`
		Temperature float64 `yaml:"temperature"`
		MaxTokens   int     `yaml:"max_tokens"`
	} `yaml:"llm"`

	Document struct {
		ChunkSize        int      `yaml:"chunk_size"`
		ChunkOverlap     int      `yaml:"chunk_overlap"`
		SupportedFormats []string `yaml:"supported_formats"`
	} `yaml:"document"`

	VectorDB struct {
		StoragePath         string  `yaml:"storage_path"`
		SimilarityThreshold float64 `yaml:"similarity_threshold"`
	} `yaml:"vector_db"`

	Logging struct {
		Level  string `yaml:"level"`
		Output string `yaml:"output"`
	} `yaml:"logging"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// GetDefaultConfig returns a default configuration
func GetDefaultConfig() *Config {
	return &Config{
		Server: struct {
			Port int    `yaml:"port"`
			Host string `yaml:"host"`
		}{
			Port: 8080,
			Host: "localhost",
		},
		Embedding: struct {
			URL       string `yaml:"url"`
			Model     string `yaml:"model"`
			BatchSize int    `yaml:"batch_size"`
		}{
			URL:       "http://localhost:8000",
			Model:     "all-MiniLM-L6-v2",
			BatchSize: 32,
		},
		LLM: struct {
			URL         string  `yaml:"url"`
			Model       string  `yaml:"model"`
			Temperature float64 `yaml:"temperature"`
			MaxTokens   int     `yaml:"max_tokens"`
		}{
			URL:         "http://localhost:11434",
			Model:       "llama2",
			Temperature: 0.7,
			MaxTokens:   512,
		},
		Document: struct {
			ChunkSize        int      `yaml:"chunk_size"`
			ChunkOverlap     int      `yaml:"chunk_overlap"`
			SupportedFormats []string `yaml:"supported_formats"`
		}{
			ChunkSize:        512,
			ChunkOverlap:     50,
			SupportedFormats: []string{"txt", "pdf", "docx"},
		},
		VectorDB: struct {
			StoragePath         string  `yaml:"storage_path"`
			SimilarityThreshold float64 `yaml:"similarity_threshold"`
		}{
			StoragePath:         "./data/vectors",
			SimilarityThreshold: 0.7,
		},
		Logging: struct {
			Level  string `yaml:"level"`
			Output string `yaml:"output"`
		}{
			Level:  "info",
			Output: "stdout",
		},
	}
}