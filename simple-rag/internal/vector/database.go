package vector

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"simple-rag/pkg/types"
)

// Database represents a simple vector database using JSON storage
type Database struct {
	storagePath string
	mu          sync.RWMutex
	chunks      map[string]*types.DocumentChunk
	documents   map[string]*types.Document
}

// NewDatabase creates a new vector database
func NewDatabase(storagePath string) *Database {
	return &Database{
		storagePath: storagePath,
		chunks:      make(map[string]*types.DocumentChunk),
		documents:   make(map[string]*types.Document),
	}
}

// Initialize creates necessary directories and loads existing data
func (db *Database) Initialize() error {
	if err := os.MkdirAll(db.storagePath, 0755); err != nil {
		return fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Load existing data
	if err := db.loadData(); err != nil {
		return fmt.Errorf("failed to load existing data: %w", err)
	}

	return nil
}

// StoreDocument stores a document in the database
func (db *Database) StoreDocument(doc *types.Document) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.documents[doc.ID] = doc
	return db.saveDocuments()
}

// StoreChunk stores a document chunk with its embedding
func (db *Database) StoreChunk(chunk *types.DocumentChunk) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.chunks[chunk.ID] = chunk
	return db.saveChunks()
}

// Search performs similarity search and returns top k results
func (db *Database) Search(queryEmbedding []float64, k int, threshold float64) ([]*types.SearchResult, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var results []*types.SearchResult

	for _, chunk := range db.chunks {
		if chunk.Embedding == nil {
			continue
		}

		similarity := cosineSimilarity(queryEmbedding, chunk.Embedding)
		if similarity >= threshold {
			document := db.documents[chunk.DocumentID]
			result := &types.SearchResult{
				Chunk:      chunk,
				Document:   document,
				Similarity: similarity,
			}
			results = append(results, result)
		}
	}

	// Sort by similarity (descending)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	// Return top k results
	if len(results) > k {
		results = results[:k]
	}

	return results, nil
}

// GetDocument retrieves a document by ID
func (db *Database) GetDocument(id string) (*types.Document, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	doc, exists := db.documents[id]
	if !exists {
		return nil, fmt.Errorf("document not found: %s", id)
	}

	return doc, nil
}

// GetChunk retrieves a chunk by ID
func (db *Database) GetChunk(id string) (*types.DocumentChunk, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	chunk, exists := db.chunks[id]
	if !exists {
		return nil, fmt.Errorf("chunk not found: %s", id)
	}

	return chunk, nil
}

// ListDocuments returns all documents
func (db *Database) ListDocuments() []*types.Document {
	db.mu.RLock()
	defer db.mu.RUnlock()

	documents := make([]*types.Document, 0, len(db.documents))
	for _, doc := range db.documents {
		documents = append(documents, doc)
	}

	return documents
}

// DeleteDocument removes a document and its chunks
func (db *Database) DeleteDocument(docID string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Remove document
	delete(db.documents, docID)

	// Remove associated chunks
	for chunkID, chunk := range db.chunks {
		if chunk.DocumentID == docID {
			delete(db.chunks, chunkID)
		}
	}

	// Save changes
	if err := db.saveDocuments(); err != nil {
		return err
	}
	return db.saveChunks()
}

// loadData loads existing data from storage
func (db *Database) loadData() error {
	// Load documents
	docsPath := filepath.Join(db.storagePath, "documents.json")
	if _, err := os.Stat(docsPath); err == nil {
		data, err := os.ReadFile(docsPath)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, &db.documents); err != nil {
			return err
		}
	}

	// Load chunks
	chunksPath := filepath.Join(db.storagePath, "chunks.json")
	if _, err := os.Stat(chunksPath); err == nil {
		data, err := os.ReadFile(chunksPath)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, &db.chunks); err != nil {
			return err
		}
	}

	return nil
}

// saveDocuments saves documents to storage
func (db *Database) saveDocuments() error {
	data, err := json.MarshalIndent(db.documents, "", "  ")
	if err != nil {
		return err
	}

	docsPath := filepath.Join(db.storagePath, "documents.json")
	return os.WriteFile(docsPath, data, 0644)
}

// saveChunks saves chunks to storage
func (db *Database) saveChunks() error {
	data, err := json.MarshalIndent(db.chunks, "", "  ")
	if err != nil {
		return err
	}

	chunksPath := filepath.Join(db.storagePath, "chunks.json")
	return os.WriteFile(chunksPath, data, 0644)
}

// cosineSimilarity calculates cosine similarity between two vectors
func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0.0
	}

	var dotProduct, normA, normB float64
	for i := range len(a) {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0.0 || normB == 0.0 {
		return 0.0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}
