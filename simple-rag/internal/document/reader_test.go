package document

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"simple-rag/pkg/types"
)

func TestNewReader(t *testing.T) {
	reader := NewReader(512, 50)
	if reader == nil {
		t.Fatal("NewReader returned nil")
	}
	if reader.chunkSize != 512 {
		t.Errorf("Expected chunk size 512, got %d", reader.chunkSize)
	}
	if reader.chunkOverlap != 50 {
		t.Errorf("Expected chunk overlap 50, got %d", reader.chunkOverlap)
	}
}

func TestReadDocument(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "document_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	reader := NewReader(512, 50)

	t.Run("ValidFile", func(t *testing.T) {
		// Create test file
		testContent := "This is a test document with some content."
		testFilePath := filepath.Join(tempDir, "test.txt")
		err := os.WriteFile(testFilePath, []byte(testContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		doc, err := reader.ReadDocument(testFilePath)
		if err != nil {
			t.Fatalf("ReadDocument failed: %v", err)
		}

		if doc == nil {
			t.Fatal("ReadDocument returned nil document")
		}

		if doc.Title != "test.txt" {
			t.Errorf("Expected title 'test.txt', got %s", doc.Title)
		}

		if doc.Content != testContent {
			t.Errorf("Expected content '%s', got '%s'", testContent, doc.Content)
		}

		if doc.FilePath != testFilePath {
			t.Errorf("Expected file path '%s', got '%s'", testFilePath, doc.FilePath)
		}

		if doc.FileType != "txt" {
			t.Errorf("Expected file type 'txt', got '%s'", doc.FileType)
		}

		if doc.FileSize != int64(len(testContent)) {
			t.Errorf("Expected file size %d, got %d", len(testContent), doc.FileSize)
		}

		if doc.Hash == "" {
			t.Error("Expected non-empty hash")
		}

		if doc.ID == "" {
			t.Error("Expected non-empty document ID")
		}

		if doc.Metadata == nil {
			t.Error("Expected non-nil metadata")
		}

		if doc.Metadata["original_name"] != "test.txt" {
			t.Errorf("Expected metadata original_name 'test.txt', got '%s'", doc.Metadata["original_name"])
		}

		if doc.Metadata["file_extension"] != "txt" {
			t.Errorf("Expected metadata file_extension 'txt', got '%s'", doc.Metadata["file_extension"])
		}

		if doc.CreatedAt.IsZero() {
			t.Error("Expected non-zero CreatedAt")
		}

		if doc.ProcessedAt.IsZero() {
			t.Error("Expected non-zero ProcessedAt")
		}
	})

	t.Run("FileWithoutExtension", func(t *testing.T) {
		testContent := "File without extension"
		testFilePath := filepath.Join(tempDir, "noext")
		err := os.WriteFile(testFilePath, []byte(testContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		doc, err := reader.ReadDocument(testFilePath)
		if err != nil {
			t.Fatalf("ReadDocument failed: %v", err)
		}

		if doc.FileType != "" {
			t.Errorf("Expected empty file type, got '%s'", doc.FileType)
		}
	})

	t.Run("NonExistentFile", func(t *testing.T) {
		_, err := reader.ReadDocument("/non/existent/file.txt")
		if err == nil {
			t.Error("Expected error for non-existent file, got nil")
		}
	})

	t.Run("EmptyFile", func(t *testing.T) {
		testFilePath := filepath.Join(tempDir, "empty.txt")
		err := os.WriteFile(testFilePath, []byte(""), 0644)
		if err != nil {
			t.Fatalf("Failed to create empty test file: %v", err)
		}

		doc, err := reader.ReadDocument(testFilePath)
		if err != nil {
			t.Fatalf("ReadDocument failed for empty file: %v", err)
		}

		if doc.Content != "" {
			t.Errorf("Expected empty content, got '%s'", doc.Content)
		}

		if doc.FileSize != 0 {
			t.Errorf("Expected file size 0, got %d", doc.FileSize)
		}
	})
}

func TestChunkDocument(t *testing.T) {
	reader := NewReader(20, 5) // Small chunks for testing

	t.Run("ValidDocument", func(t *testing.T) {
		doc := &types.Document{
			ID:      "test_doc",
			Content: "This is a test document with multiple sentences. It should be split into chunks properly.",
		}

		chunks, err := reader.ChunkDocument(doc)
		if err != nil {
			t.Fatalf("ChunkDocument failed: %v", err)
		}

		if len(chunks) == 0 {
			t.Fatal("Expected at least one chunk, got none")
		}

		// Verify chunk properties
		for i, chunk := range chunks {
			if chunk.ID == "" {
				t.Errorf("Chunk %d has empty ID", i)
			}

			if chunk.DocumentID != doc.ID {
				t.Errorf("Chunk %d has wrong document ID: expected %s, got %s", i, doc.ID, chunk.DocumentID)
			}

			if chunk.ChunkIndex != i {
				t.Errorf("Chunk %d has wrong index: expected %d, got %d", i, i, chunk.ChunkIndex)
			}

			if chunk.Content == "" {
				t.Errorf("Chunk %d has empty content", i)
			}

			if len(chunk.Content) > reader.chunkSize {
				t.Errorf("Chunk %d content too long: %d > %d", i, len(chunk.Content), reader.chunkSize)
			}

			if chunk.StartPos < 0 {
				t.Errorf("Chunk %d has negative start position: %d", i, chunk.StartPos)
			}

			if chunk.EndPos <= chunk.StartPos {
				t.Errorf("Chunk %d has invalid end position: start=%d, end=%d", i, chunk.StartPos, chunk.EndPos)
			}

			if chunk.CreatedAt.IsZero() {
				t.Errorf("Chunk %d has zero CreatedAt", i)
			}
		}

		// Verify that chunks have reasonable content coverage
		// Note: Due to chunking with overlap and word boundaries, 
		// total chunk content may not exactly match original
		totalLength := 0
		for _, chunk := range chunks {
			totalLength += len(chunk.Content)
		}
		
		// Should have some reasonable coverage
		if totalLength < len(doc.Content)/2 {
			t.Errorf("Chunks seem to have insufficient content coverage: %d chars from %d original", totalLength, len(doc.Content))
		}
	})

	t.Run("EmptyDocument", func(t *testing.T) {
		doc := &types.Document{
			ID:      "empty_doc",
			Content: "",
		}

		_, err := reader.ChunkDocument(doc)
		if err == nil {
			t.Error("Expected error for empty document, got nil")
		}
	})

	t.Run("SmallDocument", func(t *testing.T) {
		doc := &types.Document{
			ID:      "small_doc",
			Content: "Short",
		}

		chunks, err := reader.ChunkDocument(doc)
		if err != nil {
			t.Fatalf("ChunkDocument failed: %v", err)
		}

		// With chunk size 20 and overlap 5, even "Short" should result in 1 chunk
		if len(chunks) < 1 {
			t.Errorf("Expected at least 1 chunk for small document, got %d", len(chunks))
		}

		// First chunk should contain "Short"
		if !strings.Contains(chunks[0].Content, "Short") {
			t.Errorf("Expected first chunk to contain 'Short', got '%s'", chunks[0].Content)
		}
	})

	t.Run("WhitespaceOnlyDocument", func(t *testing.T) {
		doc := &types.Document{
			ID:      "whitespace_doc",
			Content: "   \n\t   \n   ",
		}

		chunks, err := reader.ChunkDocument(doc)
		if err != nil {
			t.Fatalf("ChunkDocument failed: %v", err)
		}

		// Should result in no chunks since content is only whitespace
		if len(chunks) != 0 {
			t.Errorf("Expected 0 chunks for whitespace-only document, got %d", len(chunks))
		}
	})
}

func TestGenerateDocumentID(t *testing.T) {
	id1 := generateDocumentID("/path/to/file1.txt", "hash1")
	id2 := generateDocumentID("/path/to/file2.txt", "hash2")
	id3 := generateDocumentID("/path/to/file1.txt", "hash1") // Same as id1

	if id1 == "" {
		t.Error("Generated document ID is empty")
	}

	if id1 == id2 {
		t.Error("Different files should have different IDs")
	}

	if id1 != id3 {
		t.Error("Same file and hash should generate same ID")
	}

	if len(id1) != 16 {
		t.Errorf("Expected document ID length 16, got %d", len(id1))
	}
}

func TestGenerateChunkID(t *testing.T) {
	docID := "test_doc_123"
	id1 := generateChunkID(docID, 0)
	id2 := generateChunkID(docID, 1)
	id3 := generateChunkID(docID, 0) // Same as id1

	if id1 == "" {
		t.Error("Generated chunk ID is empty")
	}

	if id1 == id2 {
		t.Error("Different chunk indices should have different IDs")
	}

	if id1 != id3 {
		t.Error("Same document ID and chunk index should generate same ID")
	}

	expectedID := "test_doc_123_chunk_0"
	if id1 != expectedID {
		t.Errorf("Expected chunk ID '%s', got '%s'", expectedID, id1)
	}
}

func TestIsSupported(t *testing.T) {
	tests := []struct {
		filePath string
		expected bool
	}{
		{"test.txt", true},
		{"test.md", true},
		{"test.text", true},
		{"test.TXT", true}, // Case insensitive
		{"test.MD", true},
		{"test.pdf", false},
		{"test.docx", false},
		{"test.html", false},
		{"test", false}, // No extension
		{"", false},     // Empty path
	}

	for _, test := range tests {
		result := IsSupported(test.filePath)
		if result != test.expected {
			t.Errorf("IsSupported(%s): expected %v, got %v", test.filePath, test.expected, result)
		}
	}
}

func TestDocumentHashConsistency(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "hash_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	reader := NewReader(512, 50)

	// Create test file
	testContent := "This is a test for hash consistency."
	testFilePath := filepath.Join(tempDir, "hash_test.txt")
	err = os.WriteFile(testFilePath, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Read document twice
	doc1, err := reader.ReadDocument(testFilePath)
	if err != nil {
		t.Fatalf("First ReadDocument failed: %v", err)
	}

	// Wait a bit to ensure different ProcessedAt times
	time.Sleep(time.Millisecond * 10)

	doc2, err := reader.ReadDocument(testFilePath)
	if err != nil {
		t.Fatalf("Second ReadDocument failed: %v", err)
	}

	// Hash should be the same
	if doc1.Hash != doc2.Hash {
		t.Errorf("Hash should be consistent: got %s and %s", doc1.Hash, doc2.Hash)
	}

	// ID should be the same (based on path and hash)
	if doc1.ID != doc2.ID {
		t.Errorf("Document ID should be consistent: got %s and %s", doc1.ID, doc2.ID)
	}

	// ProcessedAt should be different
	if doc1.ProcessedAt.Equal(doc2.ProcessedAt) {
		t.Error("ProcessedAt should be different for multiple reads")
	}
}