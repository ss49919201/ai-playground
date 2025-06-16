package types

import (
	"time"
)

// Document represents a processed document with its content and metadata
type Document struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Content     string            `json:"content"`
	FilePath    string            `json:"file_path"`
	FileType    string            `json:"file_type"`
	FileSize    int64             `json:"file_size"`
	Hash        string            `json:"hash"`
	Metadata    map[string]string `json:"metadata"`
	CreatedAt   time.Time         `json:"created_at"`
	ProcessedAt time.Time         `json:"processed_at"`
}

// DocumentChunk represents a chunk of a document with embedding
type DocumentChunk struct {
	ID         string    `json:"id"`
	DocumentID string    `json:"document_id"`
	ChunkIndex int       `json:"chunk_index"`
	Content    string    `json:"content"`
	StartPos   int       `json:"start_pos"`
	EndPos     int       `json:"end_pos"`
	Embedding  []float64 `json:"embedding"`
	CreatedAt  time.Time `json:"created_at"`
}

// SearchResult represents a search result with similarity score
type SearchResult struct {
	Chunk      *DocumentChunk `json:"chunk"`
	Document   *Document      `json:"document"`
	Similarity float64        `json:"similarity"`
}

// RAGResponse represents the response from the RAG system
type RAGResponse struct {
	Query       string          `json:"query"`
	Answer      string          `json:"answer"`
	Sources     []*SearchResult `json:"sources"`
	ProcessTime time.Duration   `json:"process_time"`
	CreatedAt   time.Time       `json:"created_at"`
}

// EmbeddingRequest represents a request to the embedding service
type EmbeddingRequest struct {
	Texts []string `json:"texts"`
	Model string   `json:"model,omitempty"`
}

// EmbeddingResponse represents a response from the embedding service
type EmbeddingResponse struct {
	Embeddings [][]float64 `json:"embeddings"`
	Model      string      `json:"model"`
}

// LLMRequest represents a request to the LLM service
type LLMRequest struct {
	Prompt      string  `json:"prompt"`
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature,omitempty"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
}

// LLMResponse represents a response from the LLM service
type LLMResponse struct {
	Response string `json:"response"`
	Model    string `json:"model"`
}