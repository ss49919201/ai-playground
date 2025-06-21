package llm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"simple-rag/pkg/types"
)

func TestNewClient(t *testing.T) {
	client := NewClient("http://test:11434", "test-model", 0.7, 512)
	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	if client.baseURL != "http://test:11434" {
		t.Errorf("Expected baseURL 'http://test:11434', got %s", client.baseURL)
	}
	if client.model != "test-model" {
		t.Errorf("Expected model 'test-model', got %s", client.model)
	}
	if client.temperature != 0.7 {
		t.Errorf("Expected temperature 0.7, got %f", client.temperature)
	}
	if client.maxTokens != 512 {
		t.Errorf("Expected maxTokens 512, got %d", client.maxTokens)
	}
	if client.httpClient == nil {
		t.Error("httpClient should be initialized")
	}
	if client.httpClient.Timeout != 60*time.Second {
		t.Errorf("Expected timeout 60s, got %v", client.httpClient.Timeout)
	}
}

func TestGenerate(t *testing.T) {
	t.Run("SuccessfulRequest", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/generate" {
				t.Errorf("Expected path '/api/generate', got %s", r.URL.Path)
			}
			if r.Method != "POST" {
				t.Errorf("Expected POST method, got %s", r.Method)
			}
			if r.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Expected Content-Type 'application/json', got %s", r.Header.Get("Content-Type"))
			}

			// Parse request
			var req OllamaRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("Failed to decode request: %v", err)
			}

			if req.Model != "test-model" {
				t.Errorf("Expected model 'test-model', got %s", req.Model)
			}
			if req.Prompt != "Test prompt" {
				t.Errorf("Expected prompt 'Test prompt', got %s", req.Prompt)
			}
			if req.Stream != false {
				t.Errorf("Expected stream false, got %v", req.Stream)
			}

			// Check options - JSON unmarshaling may result in different numeric types
			if temp, ok := req.Options["temperature"]; !ok {
				t.Error("temperature option missing")
			} else {
				// Handle both int and float64 from JSON
				var tempVal float64
				switch v := temp.(type) {
				case float64:
					tempVal = v
				case int:
					tempVal = float64(v)
				default:
					t.Errorf("Expected temperature to be numeric, got %T", temp)
				}
				if tempVal != 0.7 {
					t.Errorf("Expected temperature 0.7, got %f", tempVal)
				}
			}
			
			if tokens, ok := req.Options["num_predict"]; !ok {
				t.Error("num_predict option missing")
			} else {
				// Handle both int and float64 from JSON
				var tokensVal int
				switch v := tokens.(type) {
				case float64:
					tokensVal = int(v)
				case int:
					tokensVal = v
				default:
					t.Errorf("Expected num_predict to be numeric, got %T", tokens)
				}
				if tokensVal != 512 {
					t.Errorf("Expected num_predict 512, got %d", tokensVal)
				}
			}

			// Send response
			resp := OllamaResponse{
				Response: "This is a test response",
				Done:     true,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := NewClient(server.URL, "test-model", 0.7, 512)
		response, err := client.Generate("Test prompt")
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}

		if response != "This is a test response" {
			t.Errorf("Expected response 'This is a test response', got %s", response)
		}
	})

	t.Run("ServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewClient(server.URL, "test-model", 0.7, 512)
		_, err := client.Generate("Test prompt")
		if err == nil {
			t.Error("Expected error for server error, got nil")
		}
		if !strings.Contains(err.Error(), "status 500") {
			t.Errorf("Expected status 500 error, got %s", err.Error())
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("invalid json"))
		}))
		defer server.Close()

		client := NewClient(server.URL, "test-model", 0.7, 512)
		_, err := client.Generate("Test prompt")
		if err == nil {
			t.Error("Expected error for invalid JSON, got nil")
		}
		if !strings.Contains(err.Error(), "failed to decode response") {
			t.Errorf("Expected decode error, got %s", err.Error())
		}
	})

	t.Run("UnreachableServer", func(t *testing.T) {
		client := NewClient("http://localhost:99999", "test-model", 0.7, 512)
		_, err := client.Generate("Test prompt")
		if err == nil {
			t.Error("Expected error for unreachable server, got nil")
		}
		if !strings.Contains(err.Error(), "failed to make request") {
			t.Errorf("Expected connection error, got %s", err.Error())
		}
	})

	t.Run("WhitespaceHandling", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := OllamaResponse{
				Response: "  \n  Response with whitespace  \n  ",
				Done:     true,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := NewClient(server.URL, "test-model", 0.7, 512)
		response, err := client.Generate("Test prompt")
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}

		if response != "Response with whitespace" {
			t.Errorf("Expected trimmed response 'Response with whitespace', got '%s'", response)
		}
	})
}

func TestGenerateWithContext(t *testing.T) {
	t.Run("SuccessfulRAGQuery", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Parse request to verify RAG prompt structure
			var req OllamaRequest
			json.NewDecoder(r.Body).Decode(&req)

			// Verify that prompt contains context and question
			if !strings.Contains(req.Prompt, "Context:") {
				t.Error("Prompt should contain 'Context:'")
			}
			if !strings.Contains(req.Prompt, "Question:") {
				t.Error("Prompt should contain 'Question:'")
			}
			if !strings.Contains(req.Prompt, "What is Go?") {
				t.Error("Prompt should contain the query")
			}
			if !strings.Contains(req.Prompt, "Go is a programming language") {
				t.Error("Prompt should contain the context from search results")
			}

			resp := OllamaResponse{
				Response: "Go is a programming language developed by Google.",
				Done:     true,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := NewClient(server.URL, "test-model", 0.7, 512)

		// Create mock search results
		results := []*types.SearchResult{
			{
				Chunk: &types.DocumentChunk{
					Content: "Go is a programming language developed by Google. It's designed for simplicity and efficiency.",
				},
				Similarity: 0.95,
			},
			{
				Chunk: &types.DocumentChunk{
					Content: "Go has strong concurrency support with goroutines and channels.",
				},
				Similarity: 0.85,
			},
		}

		response, err := client.GenerateWithContext("What is Go?", results)
		if err != nil {
			t.Fatalf("GenerateWithContext failed: %v", err)
		}

		if response == nil {
			t.Fatal("Response should not be nil")
		}

		if response.Query != "What is Go?" {
			t.Errorf("Expected query 'What is Go?', got %s", response.Query)
		}

		if response.Answer != "Go is a programming language developed by Google." {
			t.Errorf("Expected answer 'Go is a programming language developed by Google.', got %s", response.Answer)
		}

		if len(response.Sources) != 2 {
			t.Errorf("Expected 2 sources, got %d", len(response.Sources))
		}

		if response.ProcessTime <= 0 {
			t.Error("ProcessTime should be positive")
		}

		if response.CreatedAt.IsZero() {
			t.Error("CreatedAt should not be zero")
		}
	})

	t.Run("EmptySearchResults", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := OllamaResponse{
				Response: "I don't have enough information to answer your question.",
				Done:     true,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := NewClient(server.URL, "test-model", 0.7, 512)
		response, err := client.GenerateWithContext("What is quantum computing?", []*types.SearchResult{})
		if err != nil {
			t.Fatalf("GenerateWithContext with empty results failed: %v", err)
		}

		if len(response.Sources) != 0 {
			t.Errorf("Expected 0 sources, got %d", len(response.Sources))
		}
	})

	t.Run("LLMError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewClient(server.URL, "test-model", 0.7, 512)
		results := []*types.SearchResult{
			{
				Chunk: &types.DocumentChunk{Content: "test content"},
				Similarity: 0.9,
			},
		}

		_, err := client.GenerateWithContext("Test query", results)
		if err == nil {
			t.Error("Expected error for LLM failure, got nil")
		}
		if !strings.Contains(err.Error(), "failed to generate response") {
			t.Errorf("Expected generation error, got %s", err.Error())
		}
	})
}

func TestBuildRAGPrompt(t *testing.T) {
	client := NewClient("http://test:11434", "test-model", 0.7, 512)

	query := "What is machine learning?"
	context := "Context 1: Machine learning is a subset of AI.\nContext 2: It involves training algorithms on data."

	prompt := client.buildRAGPrompt(query, context)

	if !strings.Contains(prompt, "Context:") {
		t.Error("Prompt should contain 'Context:'")
	}
	if !strings.Contains(prompt, "Question:") {
		t.Error("Prompt should contain 'Question:'")
	}
	if !strings.Contains(prompt, "Answer:") {
		t.Error("Prompt should contain 'Answer:'")
	}
	if !strings.Contains(prompt, query) {
		t.Error("Prompt should contain the query")
	}
	if !strings.Contains(prompt, context) {
		t.Error("Prompt should contain the context")
	}
	if !strings.Contains(prompt, "helpful assistant") {
		t.Error("Prompt should contain instruction about being a helpful assistant")
	}
}

func TestHealth(t *testing.T) {
	t.Run("HealthyService", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/tags" {
				t.Errorf("Expected path '/api/tags', got %s", r.URL.Path)
			}
			if r.Method != "GET" {
				t.Errorf("Expected GET method, got %s", r.Method)
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := NewClient(server.URL, "test-model", 0.7, 512)
		err := client.Health()
		if err != nil {
			t.Errorf("Health check should pass, got error: %v", err)
		}
	})

	t.Run("UnhealthyService", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		defer server.Close()

		client := NewClient(server.URL, "test-model", 0.7, 512)
		err := client.Health()
		if err == nil {
			t.Error("Expected error for unhealthy service, got nil")
		}
		if !strings.Contains(err.Error(), "health check failed") {
			t.Errorf("Expected health check error, got %s", err.Error())
		}
	})

	t.Run("UnreachableService", func(t *testing.T) {
		client := NewClient("http://localhost:99999", "test-model", 0.7, 512)
		err := client.Health()
		if err == nil {
			t.Error("Expected error for unreachable service, got nil")
		}
		if !strings.Contains(err.Error(), "not reachable") {
			t.Errorf("Expected unreachable error, got %s", err.Error())
		}
	})
}

func TestListModels(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/tags" {
				t.Errorf("Expected path '/api/tags', got %s", r.URL.Path)
			}

			response := struct {
				Models []struct {
					Name string `json:"name"`
				} `json:"models"`
			}{
				Models: []struct {
					Name string `json:"name"`
				}{
					{Name: "llama2"},
					{Name: "codellama"},
					{Name: "mistral"},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewClient(server.URL, "test-model", 0.7, 512)
		models, err := client.ListModels()
		if err != nil {
			t.Fatalf("ListModels failed: %v", err)
		}

		if len(models) != 3 {
			t.Errorf("Expected 3 models, got %d", len(models))
		}

		expectedModels := []string{"llama2", "codellama", "mistral"}
		for i, expected := range expectedModels {
			if models[i] != expected {
				t.Errorf("Expected model %s at index %d, got %s", expected, i, models[i])
			}
		}
	})

	t.Run("ServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewClient(server.URL, "test-model", 0.7, 512)
		_, err := client.ListModels()
		if err == nil {
			t.Error("Expected error for server error, got nil")
		}
		if !strings.Contains(err.Error(), "status 500") {
			t.Errorf("Expected status 500 error, got %s", err.Error())
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("invalid json"))
		}))
		defer server.Close()

		client := NewClient(server.URL, "test-model", 0.7, 512)
		_, err := client.ListModels()
		if err == nil {
			t.Error("Expected error for invalid JSON, got nil")
		}
		if !strings.Contains(err.Error(), "failed to decode models response") {
			t.Errorf("Expected decode error, got %s", err.Error())
		}
	})

	t.Run("UnreachableServer", func(t *testing.T) {
		client := NewClient("http://localhost:99999", "test-model", 0.7, 512)
		_, err := client.ListModels()
		if err == nil {
			t.Error("Expected error for unreachable server, got nil")
		}
		if !strings.Contains(err.Error(), "failed to list models") {
			t.Errorf("Expected connection error, got %s", err.Error())
		}
	})

	t.Run("EmptyModelsList", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := struct {
				Models []struct {
					Name string `json:"name"`
				} `json:"models"`
			}{
				Models: []struct {
					Name string `json:"name"`
				}{},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewClient(server.URL, "test-model", 0.7, 512)
		models, err := client.ListModels()
		if err != nil {
			t.Fatalf("ListModels failed: %v", err)
		}

		if len(models) != 0 {
			t.Errorf("Expected 0 models, got %d", len(models))
		}
	})
}