package vector

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"

	"simple-rag/pkg/types"
)

func TestNewDatabase(t *testing.T) {
	db := NewDatabase("/tmp/test")
	if db == nil {
		t.Fatal("NewDatabase returned nil")
	}
	if db.storagePath != "/tmp/test" {
		t.Errorf("Expected storage path '/tmp/test', got %s", db.storagePath)
	}
	if db.chunks == nil {
		t.Error("chunks map should be initialized")
	}
	if db.documents == nil {
		t.Error("documents map should be initialized")
	}
}

func TestDatabaseInitialize(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "vector_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	t.Run("NewDirectory", func(t *testing.T) {
		dbPath := filepath.Join(tempDir, "new_db")
		db := NewDatabase(dbPath)

		err := db.Initialize()
		if err != nil {
			t.Fatalf("Initialize failed: %v", err)
		}

		// Check if directory was created
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			t.Errorf("Storage directory was not created")
		}
	})

	t.Run("ExistingDirectory", func(t *testing.T) {
		dbPath := filepath.Join(tempDir, "existing_db")
		err := os.MkdirAll(dbPath, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}

		db := NewDatabase(dbPath)
		err = db.Initialize()
		if err != nil {
			t.Fatalf("Initialize failed: %v", err)
		}
	})
}

func TestStoreAndGetDocument(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "vector_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	db := NewDatabase(tempDir)
	err = db.Initialize()
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Create test document
	doc := &types.Document{
		ID:          "test_doc_1",
		Title:       "Test Document",
		Content:     "This is a test document",
		FilePath:    "/test/path.txt",
		FileType:    "txt",
		FileSize:    100,
		Hash:        "testhash",
		Metadata:    map[string]string{"key": "value"},
		CreatedAt:   time.Now(),
		ProcessedAt: time.Now(),
	}

	// Store document
	err = db.StoreDocument(doc)
	if err != nil {
		t.Fatalf("StoreDocument failed: %v", err)
	}

	// Retrieve document
	retrievedDoc, err := db.GetDocument("test_doc_1")
	if err != nil {
		t.Fatalf("GetDocument failed: %v", err)
	}

	if retrievedDoc.ID != doc.ID {
		t.Errorf("Expected ID %s, got %s", doc.ID, retrievedDoc.ID)
	}
	if retrievedDoc.Title != doc.Title {
		t.Errorf("Expected title %s, got %s", doc.Title, retrievedDoc.Title)
	}
	if retrievedDoc.Content != doc.Content {
		t.Errorf("Expected content %s, got %s", doc.Content, retrievedDoc.Content)
	}

	// Test non-existent document
	_, err = db.GetDocument("non_existent")
	if err == nil {
		t.Error("Expected error for non-existent document, got nil")
	}
}

func TestStoreAndGetChunk(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "vector_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	db := NewDatabase(tempDir)
	err = db.Initialize()
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Create test chunk
	chunk := &types.DocumentChunk{
		ID:         "test_chunk_1",
		DocumentID: "test_doc_1",
		ChunkIndex: 0,
		Content:    "This is a test chunk",
		StartPos:   0,
		EndPos:     20,
		Embedding:  []float64{0.1, 0.2, 0.3},
		CreatedAt:  time.Now(),
	}

	// Store chunk
	err = db.StoreChunk(chunk)
	if err != nil {
		t.Fatalf("StoreChunk failed: %v", err)
	}

	// Retrieve chunk
	retrievedChunk, err := db.GetChunk("test_chunk_1")
	if err != nil {
		t.Fatalf("GetChunk failed: %v", err)
	}

	if retrievedChunk.ID != chunk.ID {
		t.Errorf("Expected ID %s, got %s", chunk.ID, retrievedChunk.ID)
	}
	if retrievedChunk.DocumentID != chunk.DocumentID {
		t.Errorf("Expected DocumentID %s, got %s", chunk.DocumentID, retrievedChunk.DocumentID)
	}
	if retrievedChunk.Content != chunk.Content {
		t.Errorf("Expected content %s, got %s", chunk.Content, retrievedChunk.Content)
	}
	if len(retrievedChunk.Embedding) != len(chunk.Embedding) {
		t.Errorf("Expected embedding length %d, got %d", len(chunk.Embedding), len(retrievedChunk.Embedding))
	}

	// Test non-existent chunk
	_, err = db.GetChunk("non_existent")
	if err == nil {
		t.Error("Expected error for non-existent chunk, got nil")
	}
}

func TestSearch(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "vector_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	db := NewDatabase(tempDir)
	err = db.Initialize()
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Store test document
	doc := &types.Document{
		ID:      "test_doc_1",
		Title:   "Test Document",
		Content: "This is a test document",
	}
	err = db.StoreDocument(doc)
	if err != nil {
		t.Fatalf("StoreDocument failed: %v", err)
	}

	// Store test chunks with embeddings
	chunks := []*types.DocumentChunk{
		{
			ID:         "chunk_1",
			DocumentID: "test_doc_1",
			Content:    "First chunk",
			Embedding:  []float64{1.0, 0.0, 0.0}, // Similar to query
		},
		{
			ID:         "chunk_2",
			DocumentID: "test_doc_1",
			Content:    "Second chunk",
			Embedding:  []float64{0.0, 1.0, 0.0}, // Less similar to query
		},
		{
			ID:         "chunk_3",
			DocumentID: "test_doc_1",
			Content:    "Third chunk",
			Embedding:  []float64{0.5, 0.5, 0.0}, // Moderate similarity
		},
		{
			ID:         "chunk_4",
			DocumentID: "test_doc_1",
			Content:    "Fourth chunk",
			// No embedding - should be skipped
		},
	}

	for _, chunk := range chunks {
		err = db.StoreChunk(chunk)
		if err != nil {
			t.Fatalf("StoreChunk failed: %v", err)
		}
	}

	// Test search
	queryEmbedding := []float64{0.9, 0.1, 0.0} // Most similar to chunk_1
	results, err := db.Search(queryEmbedding, 10, 0.0)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	// Should return 3 results (chunk_4 has no embedding)
	if len(results) != 3 {
		t.Errorf("Expected 3 search results, got %d", len(results))
	}

	// Results should be sorted by similarity (descending)
	for i := 0; i < len(results)-1; i++ {
		if results[i].Similarity < results[i+1].Similarity {
			t.Errorf("Results not sorted by similarity: %f < %f", results[i].Similarity, results[i+1].Similarity)
		}
	}

	// First result should be chunk_1 (highest similarity)
	if results[0].Chunk.ID != "chunk_1" {
		t.Errorf("Expected first result to be chunk_1, got %s", results[0].Chunk.ID)
	}

	// Test with threshold
	results, err = db.Search(queryEmbedding, 10, 0.8)
	if err != nil {
		t.Fatalf("Search with threshold failed: %v", err)
	}

	// Only chunk_1 should meet the threshold
	if len(results) != 1 {
		t.Errorf("Expected 1 result with high threshold, got %d", len(results))
	}

	// Test with k limit
	results, err = db.Search(queryEmbedding, 2, 0.0)
	if err != nil {
		t.Fatalf("Search with k limit failed: %v", err)
	}

	// Should return only 2 results
	if len(results) != 2 {
		t.Errorf("Expected 2 results with k=2, got %d", len(results))
	}
}

func TestListDocuments(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "vector_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	db := NewDatabase(tempDir)
	err = db.Initialize()
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Initially should be empty
	docs := db.ListDocuments()
	if len(docs) != 0 {
		t.Errorf("Expected 0 documents initially, got %d", len(docs))
	}

	// Add documents
	for i := 0; i < 3; i++ {
		doc := &types.Document{
			ID:      fmt.Sprintf("doc_%d", i),
			Title:   fmt.Sprintf("Document %d", i),
			Content: fmt.Sprintf("Content of document %d", i),
		}
		err = db.StoreDocument(doc)
		if err != nil {
			t.Fatalf("StoreDocument failed: %v", err)
		}
	}

	// Should have 3 documents
	docs = db.ListDocuments()
	if len(docs) != 3 {
		t.Errorf("Expected 3 documents, got %d", len(docs))
	}
}

func TestDeleteDocument(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "vector_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	db := NewDatabase(tempDir)
	err = db.Initialize()
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Store document and chunks
	doc := &types.Document{
		ID:      "test_doc_1",
		Title:   "Test Document",
		Content: "This is a test document",
	}
	err = db.StoreDocument(doc)
	if err != nil {
		t.Fatalf("StoreDocument failed: %v", err)
	}

	chunks := []*types.DocumentChunk{
		{
			ID:         "chunk_1",
			DocumentID: "test_doc_1",
			Content:    "First chunk",
		},
		{
			ID:         "chunk_2",
			DocumentID: "test_doc_1",
			Content:    "Second chunk",
		},
	}

	for _, chunk := range chunks {
		err = db.StoreChunk(chunk)
		if err != nil {
			t.Fatalf("StoreChunk failed: %v", err)
		}
	}

	// Verify data exists
	_, err = db.GetDocument("test_doc_1")
	if err != nil {
		t.Fatalf("Document should exist before deletion: %v", err)
	}
	_, err = db.GetChunk("chunk_1")
	if err != nil {
		t.Fatalf("Chunk should exist before deletion: %v", err)
	}

	// Delete document
	err = db.DeleteDocument("test_doc_1")
	if err != nil {
		t.Fatalf("DeleteDocument failed: %v", err)
	}

	// Verify document is deleted
	_, err = db.GetDocument("test_doc_1")
	if err == nil {
		t.Error("Document should not exist after deletion")
	}

	// Verify chunks are deleted
	_, err = db.GetChunk("chunk_1")
	if err == nil {
		t.Error("Chunk should not exist after document deletion")
	}
	_, err = db.GetChunk("chunk_2")
	if err == nil {
		t.Error("Chunk should not exist after document deletion")
	}
}

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a        []float64
		b        []float64
		expected float64
	}{
		{
			name:     "Identical vectors",
			a:        []float64{1.0, 0.0, 0.0},
			b:        []float64{1.0, 0.0, 0.0},
			expected: 1.0,
		},
		{
			name:     "Orthogonal vectors",
			a:        []float64{1.0, 0.0, 0.0},
			b:        []float64{0.0, 1.0, 0.0},
			expected: 0.0,
		},
		{
			name:     "Opposite vectors",
			a:        []float64{1.0, 0.0, 0.0},
			b:        []float64{-1.0, 0.0, 0.0},
			expected: -1.0,
		},
		{
			name:     "Different lengths",
			a:        []float64{1.0, 0.0},
			b:        []float64{1.0, 0.0, 0.0},
			expected: 0.0,
		},
		{
			name:     "Zero vector",
			a:        []float64{0.0, 0.0, 0.0},
			b:        []float64{1.0, 0.0, 0.0},
			expected: 0.0,
		},
		{
			name:     "Similar vectors",
			a:        []float64{0.6, 0.8, 0.0},
			b:        []float64{0.8, 0.6, 0.0},
			expected: 0.96, // cos(angle) between these vectors
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := cosineSimilarity(test.a, test.b)
			if math.Abs(result-test.expected) > 1e-6 {
				t.Errorf("Expected similarity %f, got %f", test.expected, result)
			}
		})
	}
}

func TestDataPersistence(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "vector_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create and populate database
	db1 := NewDatabase(tempDir)
	err = db1.Initialize()
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	doc := &types.Document{
		ID:      "persistent_doc",
		Title:   "Persistent Document",
		Content: "This document should persist",
	}
	err = db1.StoreDocument(doc)
	if err != nil {
		t.Fatalf("StoreDocument failed: %v", err)
	}

	chunk := &types.DocumentChunk{
		ID:         "persistent_chunk",
		DocumentID: "persistent_doc",
		Content:    "This chunk should persist",
		Embedding:  []float64{0.1, 0.2, 0.3},
	}
	err = db1.StoreChunk(chunk)
	if err != nil {
		t.Fatalf("StoreChunk failed: %v", err)
	}

	// Create new database instance with same path
	db2 := NewDatabase(tempDir)
	err = db2.Initialize()
	if err != nil {
		t.Fatalf("Second Initialize failed: %v", err)
	}

	// Data should be loaded from storage
	retrievedDoc, err := db2.GetDocument("persistent_doc")
	if err != nil {
		t.Fatalf("GetDocument failed after reload: %v", err)
	}
	if retrievedDoc.Title != doc.Title {
		t.Errorf("Document not properly persisted: expected %s, got %s", doc.Title, retrievedDoc.Title)
	}

	retrievedChunk, err := db2.GetChunk("persistent_chunk")
	if err != nil {
		t.Fatalf("GetChunk failed after reload: %v", err)
	}
	if retrievedChunk.Content != chunk.Content {
		t.Errorf("Chunk not properly persisted: expected %s, got %s", chunk.Content, retrievedChunk.Content)
	}
	if len(retrievedChunk.Embedding) != len(chunk.Embedding) {
		t.Errorf("Chunk embedding not properly persisted")
	}
}