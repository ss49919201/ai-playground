package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"simple-rag/pkg/types"
)

// Client handles communication with Ollama LLM service
type Client struct {
	baseURL     string
	model       string
	temperature float64
	maxTokens   int
	httpClient  *http.Client
}

// NewClient creates a new LLM client
func NewClient(baseURL, model string, temperature float64, maxTokens int) *Client {
	return &Client{
		baseURL:     baseURL,
		model:       model,
		temperature: temperature,
		maxTokens:   maxTokens,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// OllamaRequest represents an Ollama API request
type OllamaRequest struct {
	Model       string                 `json:"model"`
	Prompt      string                 `json:"prompt"`
	Stream      bool                   `json:"stream"`
	Options     map[string]interface{} `json:"options,omitempty"`
}

// OllamaResponse represents an Ollama API response
type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// Generate generates a response using the LLM
func (c *Client) Generate(prompt string) (string, error) {
	// Prepare request
	req := OllamaRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: false,
		Options: map[string]interface{}{
			"temperature": c.temperature,
			"num_predict": c.maxTokens,
		},
	}

	// Serialize request
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make HTTP request
	url := fmt.Sprintf("%s/api/generate", c.baseURL)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM service returned status %d", resp.StatusCode)
	}

	// Parse response
	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return strings.TrimSpace(ollamaResp.Response), nil
}

// GenerateWithContext generates a response using retrieved context
func (c *Client) GenerateWithContext(query string, searchResults []*types.SearchResult) (*types.RAGResponse, error) {
	startTime := time.Now()

	// Build context from search results
	var contextParts []string
	for i, result := range searchResults {
		contextParts = append(contextParts, fmt.Sprintf("Context %d (similarity: %.3f):\n%s", 
			i+1, result.Similarity, result.Chunk.Content))
	}
	context := strings.Join(contextParts, "\n\n")

	// Build prompt
	prompt := c.buildRAGPrompt(query, context)

	// Generate response
	answer, err := c.Generate(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate response: %w", err)
	}

	// Build response
	response := &types.RAGResponse{
		Query:       query,
		Answer:      answer,
		Sources:     searchResults,
		ProcessTime: time.Since(startTime),
		CreatedAt:   time.Now(),
	}

	return response, nil
}

// buildRAGPrompt builds a prompt for RAG using retrieved context
func (c *Client) buildRAGPrompt(query, context string) string {
	template := `You are a helpful assistant that answers questions based on the provided context. 
Use only the information from the context to answer the question. If the context doesn't contain enough information to answer the question, say so.

Context:
%s

Question: %s

Answer:`

	return fmt.Sprintf(template, context, query)
}

// Health checks if the LLM service is available
func (c *Client) Health() error {
	url := fmt.Sprintf("%s/api/tags", c.baseURL)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("LLM service not reachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("LLM service health check failed: status %d", resp.StatusCode)
	}

	return nil
}

// ListModels lists available models
func (c *Client) ListModels() ([]string, error) {
	url := fmt.Sprintf("%s/api/tags", c.baseURL)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list models: status %d", resp.StatusCode)
	}

	var response struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode models response: %w", err)
	}

	models := make([]string, len(response.Models))
	for i, model := range response.Models {
		models[i] = model.Name
	}

	return models, nil
}