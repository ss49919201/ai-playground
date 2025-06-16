package vector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"simple-rag/pkg/types"
)

// EmbeddingClient handles communication with embedding service
type EmbeddingClient struct {
	baseURL    string
	model      string
	httpClient *http.Client
}

// NewEmbeddingClient creates a new embedding client
func NewEmbeddingClient(baseURL, model string) *EmbeddingClient {
	return &EmbeddingClient{
		baseURL: baseURL,
		model:   model,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetEmbeddings gets embeddings for a list of texts
func (c *EmbeddingClient) GetEmbeddings(texts []string) ([][]float64, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("no texts provided")
	}

	// Prepare request
	req := types.EmbeddingRequest{
		Texts: texts,
		Model: c.model,
	}

	// Serialize request
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make HTTP request
	url := fmt.Sprintf("%s/embeddings", c.baseURL)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("embedding service returned status %d", resp.StatusCode)
	}

	// Parse response
	var embResp types.EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(embResp.Embeddings) != len(texts) {
		return nil, fmt.Errorf("embedding count mismatch: expected %d, got %d", 
			len(texts), len(embResp.Embeddings))
	}

	return embResp.Embeddings, nil
}

// GetSingleEmbedding gets embedding for a single text
func (c *EmbeddingClient) GetSingleEmbedding(text string) ([]float64, error) {
	embeddings, err := c.GetEmbeddings([]string{text})
	if err != nil {
		return nil, err
	}
	
	if len(embeddings) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}
	
	return embeddings[0], nil
}

// ProcessChunks processes document chunks and adds embeddings
func (c *EmbeddingClient) ProcessChunks(chunks []*types.DocumentChunk) error {
	if len(chunks) == 0 {
		return nil
	}

	// Extract texts
	texts := make([]string, len(chunks))
	for i, chunk := range chunks {
		texts[i] = chunk.Content
	}

	// Get embeddings
	embeddings, err := c.GetEmbeddings(texts)
	if err != nil {
		return fmt.Errorf("failed to get embeddings: %w", err)
	}

	// Assign embeddings to chunks
	for i, chunk := range chunks {
		chunk.Embedding = embeddings[i]
	}

	return nil
}

// Health checks if the embedding service is available
func (c *EmbeddingClient) Health() error {
	url := fmt.Sprintf("%s/health", c.baseURL)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("embedding service not reachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("embedding service health check failed: status %d", resp.StatusCode)
	}

	return nil
}