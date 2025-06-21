package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetDefaultConfig(t *testing.T) {
	config := GetDefaultConfig()

	if config == nil {
		t.Fatal("GetDefaultConfig returned nil")
	}

	// Test server defaults
	if config.Server.Port != 8080 {
		t.Errorf("Expected server port 8080, got %d", config.Server.Port)
	}
	if config.Server.Host != "localhost" {
		t.Errorf("Expected server host 'localhost', got %s", config.Server.Host)
	}

	// Test embedding defaults
	if config.Embedding.URL != "http://localhost:8000" {
		t.Errorf("Expected embedding URL 'http://localhost:8000', got %s", config.Embedding.URL)
	}
	if config.Embedding.Model != "all-MiniLM-L6-v2" {
		t.Errorf("Expected embedding model 'all-MiniLM-L6-v2', got %s", config.Embedding.Model)
	}
	if config.Embedding.BatchSize != 32 {
		t.Errorf("Expected embedding batch size 32, got %d", config.Embedding.BatchSize)
	}

	// Test LLM defaults
	if config.LLM.URL != "http://localhost:11434" {
		t.Errorf("Expected LLM URL 'http://localhost:11434', got %s", config.LLM.URL)
	}
	if config.LLM.Model != "llama2" {
		t.Errorf("Expected LLM model 'llama2', got %s", config.LLM.Model)
	}
	if config.LLM.Temperature != 0.7 {
		t.Errorf("Expected LLM temperature 0.7, got %f", config.LLM.Temperature)
	}
	if config.LLM.MaxTokens != 512 {
		t.Errorf("Expected LLM max tokens 512, got %d", config.LLM.MaxTokens)
	}

	// Test document defaults
	if config.Document.ChunkSize != 512 {
		t.Errorf("Expected document chunk size 512, got %d", config.Document.ChunkSize)
	}
	if config.Document.ChunkOverlap != 50 {
		t.Errorf("Expected document chunk overlap 50, got %d", config.Document.ChunkOverlap)
	}
	expectedFormats := []string{"txt", "pdf", "docx"}
	if len(config.Document.SupportedFormats) != len(expectedFormats) {
		t.Errorf("Expected %d supported formats, got %d", len(expectedFormats), len(config.Document.SupportedFormats))
	}
	for i, format := range expectedFormats {
		if config.Document.SupportedFormats[i] != format {
			t.Errorf("Expected supported format %s at index %d, got %s", format, i, config.Document.SupportedFormats[i])
		}
	}

	// Test vector DB defaults
	if config.VectorDB.StoragePath != "./data/vectors" {
		t.Errorf("Expected vector DB storage path './data/vectors', got %s", config.VectorDB.StoragePath)
	}
	if config.VectorDB.SimilarityThreshold != 0.7 {
		t.Errorf("Expected vector DB similarity threshold 0.7, got %f", config.VectorDB.SimilarityThreshold)
	}

	// Test logging defaults
	if config.Logging.Level != "info" {
		t.Errorf("Expected logging level 'info', got %s", config.Logging.Level)
	}
	if config.Logging.Output != "stdout" {
		t.Errorf("Expected logging output 'stdout', got %s", config.Logging.Output)
	}
}

func TestLoadConfig(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test loading valid config
	t.Run("ValidConfig", func(t *testing.T) {
		configPath := filepath.Join(tempDir, "config.yaml")
		configContent := `
server:
  port: 9090
  host: "0.0.0.0"
embedding:
  url: "http://test:8000"
  model: "test-model"
  batch_size: 16
llm:
  url: "http://test:11434"
  model: "test-llm"
  temperature: 0.5
  max_tokens: 256
document:
  chunk_size: 256
  chunk_overlap: 25
  supported_formats: ["txt", "md"]
vector_db:
  storage_path: "/tmp/vectors"
  similarity_threshold: 0.8
logging:
  level: "debug"
  output: "file"
`

		err := os.WriteFile(configPath, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write test config file: %v", err)
		}

		config, err := LoadConfig(configPath)
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}

		if config.Server.Port != 9090 {
			t.Errorf("Expected server port 9090, got %d", config.Server.Port)
		}
		if config.Server.Host != "0.0.0.0" {
			t.Errorf("Expected server host '0.0.0.0', got %s", config.Server.Host)
		}
		if config.Embedding.URL != "http://test:8000" {
			t.Errorf("Expected embedding URL 'http://test:8000', got %s", config.Embedding.URL)
		}
		if config.Document.ChunkSize != 256 {
			t.Errorf("Expected document chunk size 256, got %d", config.Document.ChunkSize)
		}
		if config.VectorDB.SimilarityThreshold != 0.8 {
			t.Errorf("Expected vector DB similarity threshold 0.8, got %f", config.VectorDB.SimilarityThreshold)
		}
	})

	// Test loading non-existent file
	t.Run("NonExistentFile", func(t *testing.T) {
		_, err := LoadConfig("/non/existent/file.yaml")
		if err == nil {
			t.Error("Expected error for non-existent file, got nil")
		}
	})

	// Test loading invalid YAML
	t.Run("InvalidYAML", func(t *testing.T) {
		configPath := filepath.Join(tempDir, "invalid.yaml")
		invalidContent := `
server:
  port: not_a_number
  host: "localhost"
embedding:
  url: "http://localhost:8000"
  model: "test-model"
  batch_size: not_a_number
`

		err := os.WriteFile(configPath, []byte(invalidContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write invalid config file: %v", err)
		}

		_, err = LoadConfig(configPath)
		if err == nil {
			t.Error("Expected error for invalid YAML, got nil")
		}
	})

	// Test loading empty file
	t.Run("EmptyFile", func(t *testing.T) {
		configPath := filepath.Join(tempDir, "empty.yaml")
		err := os.WriteFile(configPath, []byte(""), 0644)
		if err != nil {
			t.Fatalf("Failed to write empty config file: %v", err)
		}

		config, err := LoadConfig(configPath)
		if err != nil {
			t.Fatalf("LoadConfig failed for empty file: %v", err)
		}

		// Empty file should result in zero values
		if config.Server.Port != 0 {
			t.Errorf("Expected server port 0 for empty config, got %d", config.Server.Port)
		}
		if config.Server.Host != "" {
			t.Errorf("Expected empty server host for empty config, got %s", config.Server.Host)
		}
	})
}