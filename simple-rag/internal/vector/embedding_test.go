package vector

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"simple-rag/pkg/types"
)

func TestNewEmbeddingClient(t *testing.T) {
	client := NewEmbeddingClient("http://test:8000", "test-model")
	if client == nil {
		t.Fatal("NewEmbeddingClient returned nil")
	}
	if client.baseURL != "http://test:8000" {
		t.Errorf("Expected baseURL 'http://test:8000', got %s", client.baseURL)
	}
	if client.model != "test-model" {
		t.Errorf("Expected model 'test-model', got %s", client.model)
	}
	if client.httpClient == nil {
		t.Error("httpClient should be initialized")
	}
	if client.httpClient.Timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", client.httpClient.Timeout)
	}
}

func TestGetEmbeddings(t *testing.T) {
	t.Run("SuccessfulRequest", func(t *testing.T) {
		// Create test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/embeddings" {
				t.Errorf("Expected path '/embeddings', got %s", r.URL.Path)
			}
			if r.Method != "POST" {
				t.Errorf("Expected POST method, got %s", r.Method)
			}
			if r.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Expected Content-Type 'application/json', got %s", r.Header.Get("Content-Type"))
			}

			// Parse request
			var req types.EmbeddingRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("Failed to decode request: %v", err)
			}

			if len(req.Texts) != 2 {
				t.Errorf("Expected 2 texts, got %d", len(req.Texts))
			}
			if req.Model != "test-model" {
				t.Errorf("Expected model 'test-model', got %s", req.Model)
			}

			// Send response
			resp := types.EmbeddingResponse{
				Embeddings: [][]float64{
					{0.1, 0.2, 0.3},
					{0.4, 0.5, 0.6},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := NewEmbeddingClient(server.URL, "test-model")
		texts := []string{"first text", "second text"}

		embeddings, err := client.GetEmbeddings(texts)
		if err != nil {
			t.Fatalf("GetEmbeddings failed: %v", err)
		}

		if len(embeddings) != 2 {
			t.Errorf("Expected 2 embeddings, got %d", len(embeddings))
		}

		if len(embeddings[0]) != 3 {
			t.Errorf("Expected embedding dimension 3, got %d", len(embeddings[0]))
		}

		expected := [][]float64{
			{0.1, 0.2, 0.3},
			{0.4, 0.5, 0.6},
		}

		for i, embedding := range embeddings {
			for j, val := range embedding {
				if val != expected[i][j] {
					t.Errorf("Expected embedding[%d][%d] = %f, got %f", i, j, expected[i][j], val)
				}
			}
		}
	})

	t.Run("EmptyTexts", func(t *testing.T) {
		client := NewEmbeddingClient("http://test:8000", "test-model")
		_, err := client.GetEmbeddings([]string{})
		if err == nil {
			t.Error("Expected error for empty texts, got nil")
		}
		if !strings.Contains(err.Error(), "no texts provided") {
			t.Errorf("Expected 'no texts provided' error, got %s", err.Error())
		}
	})

	t.Run("ServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewEmbeddingClient(server.URL, "test-model")
		_, err := client.GetEmbeddings([]string{"test"})
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

		client := NewEmbeddingClient(server.URL, "test-model")
		_, err := client.GetEmbeddings([]string{"test"})
		if err == nil {
			t.Error("Expected error for invalid JSON, got nil")
		}
		if !strings.Contains(err.Error(), "failed to decode response") {
			t.Errorf("Expected decode error, got %s", err.Error())
		}
	})

	t.Run("EmbeddingCountMismatch", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := types.EmbeddingResponse{
				Embeddings: [][]float64{
					{0.1, 0.2, 0.3}, // Only 1 embedding for 2 texts
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := NewEmbeddingClient(server.URL, "test-model")
		_, err := client.GetEmbeddings([]string{"text1", "text2"})
		if err == nil {
			t.Error("Expected error for embedding count mismatch, got nil")
		}
		if !strings.Contains(err.Error(), "embedding count mismatch") {
			t.Errorf("Expected count mismatch error, got %s", err.Error())
		}
	})

	t.Run("UnreachableServer", func(t *testing.T) {
		client := NewEmbeddingClient("http://localhost:99999", "test-model")
		_, err := client.GetEmbeddings([]string{"test"})
		if err == nil {
			t.Error("Expected error for unreachable server, got nil")
		}
		if !strings.Contains(err.Error(), "failed to make request") {
			t.Errorf("Expected connection error, got %s", err.Error())
		}
	})
}

func TestGetSingleEmbedding(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := types.EmbeddingResponse{
				Embeddings: [][]float64{
					{0.1, 0.2, 0.3},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := NewEmbeddingClient(server.URL, "test-model")
		embedding, err := client.GetSingleEmbedding("test text")
		if err != nil {
			t.Fatalf("GetSingleEmbedding failed: %v", err)
		}

		if len(embedding) != 3 {
			t.Errorf("Expected embedding dimension 3, got %d", len(embedding))
		}

		expected := []float64{0.1, 0.2, 0.3}
		for i, val := range embedding {
			if val != expected[i] {
				t.Errorf("Expected embedding[%d] = %f, got %f", i, expected[i], val)
			}
		}
	})

	t.Run("NoEmbedding", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := types.EmbeddingResponse{
				Embeddings: [][]float64{},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := NewEmbeddingClient(server.URL, "test-model")
		_, err := client.GetSingleEmbedding("test text")
		if err == nil {
			t.Error("Expected error for no embedding, got nil")
		}
		if !strings.Contains(err.Error(), "embedding count mismatch") && !strings.Contains(err.Error(), "no embedding returned") {
			t.Errorf("Expected embedding count mismatch or no embedding error, got %s", err.Error())
		}
	})
}

func TestProcessChunks(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := types.EmbeddingResponse{
				Embeddings: [][]float64{
					{0.1, 0.2, 0.3},
					{0.4, 0.5, 0.6},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := NewEmbeddingClient(server.URL, "test-model")

		chunks := []*types.DocumentChunk{
			{
				ID:      "chunk1",
				Content: "First chunk content",
			},
			{
				ID:      "chunk2",
				Content: "Second chunk content",
			},
		}

		err := client.ProcessChunks(chunks)
		if err != nil {
			t.Fatalf("ProcessChunks failed: %v", err)
		}

		// Check that embeddings were assigned
		if chunks[0].Embedding == nil {
			t.Error("First chunk should have embedding")
		}
		if chunks[1].Embedding == nil {
			t.Error("Second chunk should have embedding")
		}

		if len(chunks[0].Embedding) != 3 {
			t.Errorf("Expected embedding dimension 3, got %d", len(chunks[0].Embedding))
		}

		expected := [][]float64{
			{0.1, 0.2, 0.3},
			{0.4, 0.5, 0.6},
		}

		for i, chunk := range chunks {
			for j, val := range chunk.Embedding {
				if val != expected[i][j] {
					t.Errorf("Expected chunk[%d].Embedding[%d] = %f, got %f", i, j, expected[i][j], val)
				}
			}
		}
	})

	t.Run("EmptyChunks", func(t *testing.T) {
		client := NewEmbeddingClient("http://test:8000", "test-model")
		err := client.ProcessChunks([]*types.DocumentChunk{})
		if err != nil {
			t.Errorf("ProcessChunks should handle empty chunks gracefully, got error: %v", err)
		}
	})

	t.Run("EmbeddingError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewEmbeddingClient(server.URL, "test-model")

		chunks := []*types.DocumentChunk{
			{
				ID:      "chunk1",
				Content: "First chunk content",
			},
		}

		err := client.ProcessChunks(chunks)
		if err == nil {
			t.Error("Expected error for embedding failure, got nil")
		}
		if !strings.Contains(err.Error(), "failed to get embeddings") {
			t.Errorf("Expected embedding error, got %s", err.Error())
		}
	})
}

func TestHealth(t *testing.T) {
	t.Run("HealthyService", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/health" {
				t.Errorf("Expected path '/health', got %s", r.URL.Path)
			}
			if r.Method != "GET" {
				t.Errorf("Expected GET method, got %s", r.Method)
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := NewEmbeddingClient(server.URL, "test-model")
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

		client := NewEmbeddingClient(server.URL, "test-model")
		err := client.Health()
		if err == nil {
			t.Error("Expected error for unhealthy service, got nil")
		}
		if !strings.Contains(err.Error(), "health check failed") {
			t.Errorf("Expected health check error, got %s", err.Error())
		}
	})

	t.Run("UnreachableService", func(t *testing.T) {
		client := NewEmbeddingClient("http://localhost:99999", "test-model")
		err := client.Health()
		if err == nil {
			t.Error("Expected error for unreachable service, got nil")
		}
		if !strings.Contains(err.Error(), "not reachable") {
			t.Errorf("Expected unreachable error, got %s", err.Error())
		}
	})
}

func TestEmbeddingClientTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping timeout test in short mode")
	}

	// Create client with shorter timeout for testing
	client := NewEmbeddingClient("http://localhost:99999", "test-model")
	client.httpClient.Timeout = 100 * time.Millisecond

	_, err := client.GetEmbeddings([]string{"test"})
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}