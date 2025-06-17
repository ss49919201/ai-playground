package main

import (
	"fmt"
	"time"

	"simple-rag/internal/config"
	"simple-rag/internal/document"
	"simple-rag/internal/llm"
	"simple-rag/internal/vector"
	"simple-rag/pkg/types"
)

// RAGSystem combines all components for the RAG functionality
type RAGSystem struct {
	db              *vector.Database
	embeddingClient *vector.EmbeddingClient
	llmClient       *llm.Client
	docReader       *document.Reader
	config          *config.Config
}

// AddDocument processes and adds a document to the system
func (r *RAGSystem) AddDocument(filePath string) error {
	// Check if file type is supported
	if !document.IsSupported(filePath) {
		return fmt.Errorf("file type not supported: %s", filePath)
	}

	// Read document
	doc, err := r.docReader.ReadDocument(filePath)
	if err != nil {
		return fmt.Errorf("failed to read document: %w", err)
	}

	// Check if document already exists
	if existingDoc, _ := r.db.GetDocument(doc.ID); existingDoc != nil {
		return fmt.Errorf("document already exists: %s", doc.Title)
	}

	// Store document
	if err := r.db.StoreDocument(doc); err != nil {
		return fmt.Errorf("failed to store document: %w", err)
	}

	// Chunk document
	chunks, err := r.docReader.ChunkDocument(doc)
	if err != nil {
		return fmt.Errorf("failed to chunk document: %w", err)
	}

	// Process chunks in batches to get embeddings
	batchSize := r.config.Embedding.BatchSize
	for i := 0; i < len(chunks); i += batchSize {
		end := i + batchSize
		if end > len(chunks) {
			end = len(chunks)
		}

		batch := chunks[i:end]
		if err := r.embeddingClient.ProcessChunks(batch); err != nil {
			return fmt.Errorf("failed to process embeddings for batch: %w", err)
		}

		// Store chunks with embeddings
		for _, chunk := range batch {
			if err := r.db.StoreChunk(chunk); err != nil {
				return fmt.Errorf("failed to store chunk: %w", err)
			}
		}
	}

	return nil
}

// Query performs a RAG query and returns the response
func (r *RAGSystem) Query(query string) (*types.RAGResponse, error) {
	// Get query embedding
	queryEmbedding, err := r.embeddingClient.GetSingleEmbedding(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get query embedding: %w", err)
	}

	// Search for similar chunks
	searchResults, err := r.db.Search(queryEmbedding, 5, r.config.VectorDB.SimilarityThreshold)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	if len(searchResults) == 0 {
		return &types.RAGResponse{
			Query:       query,
			Answer:      "I don't have enough information to answer this question based on the available documents.",
			Sources:     searchResults,
			ProcessTime: 0,
			CreatedAt:   time.Now(),
		}, nil
	}

	// Generate response using LLM
	response, err := r.llmClient.GenerateWithContext(query, searchResults)
	if err != nil {
		return nil, fmt.Errorf("failed to generate response: %w", err)
	}

	return response, nil
}

// ListDocuments returns all documents in the system
func (r *RAGSystem) ListDocuments() []*types.Document {
	return r.db.ListDocuments()
}

// DeleteDocument removes a document and its chunks
func (r *RAGSystem) DeleteDocument(docID string) error {
	return r.db.DeleteDocument(docID)
}

// Health checks the health of all components
func (r *RAGSystem) Health() error {
	// Check embedding service
	if err := r.embeddingClient.Health(); err != nil {
		return fmt.Errorf("embedding service unhealthy: %w", err)
	}

	// Check LLM service
	if err := r.llmClient.Health(); err != nil {
		return fmt.Errorf("LLM service unhealthy: %w", err)
	}

	return nil
}