package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDocumentSerialization(t *testing.T) {
	// Create test document
	doc := &Document{
		ID:          "test_doc_123",
		Title:       "Test Document",
		Content:     "This is test content for the document",
		FilePath:    "/path/to/test.txt",
		FileType:    "txt",
		FileSize:    1024,
		Hash:        "abc123def456",
		Metadata:    map[string]string{"author": "test", "version": "1.0"},
		CreatedAt:   time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
		ProcessedAt: time.Date(2023, 12, 25, 10, 35, 0, 0, time.UTC),
	}

	// Test JSON serialization
	jsonData, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("Failed to marshal document: %v", err)
	}

	// Test JSON deserialization
	var deserializedDoc Document
	err = json.Unmarshal(jsonData, &deserializedDoc)
	if err != nil {
		t.Fatalf("Failed to unmarshal document: %v", err)
	}

	// Verify fields
	if deserializedDoc.ID != doc.ID {
		t.Errorf("Expected ID %s, got %s", doc.ID, deserializedDoc.ID)
	}
	if deserializedDoc.Title != doc.Title {
		t.Errorf("Expected title %s, got %s", doc.Title, deserializedDoc.Title)
	}
	if deserializedDoc.Content != doc.Content {
		t.Errorf("Expected content %s, got %s", doc.Content, deserializedDoc.Content)
	}
	if deserializedDoc.FilePath != doc.FilePath {
		t.Errorf("Expected file path %s, got %s", doc.FilePath, deserializedDoc.FilePath)
	}
	if deserializedDoc.FileType != doc.FileType {
		t.Errorf("Expected file type %s, got %s", doc.FileType, deserializedDoc.FileType)
	}
	if deserializedDoc.FileSize != doc.FileSize {
		t.Errorf("Expected file size %d, got %d", doc.FileSize, deserializedDoc.FileSize)
	}
	if deserializedDoc.Hash != doc.Hash {
		t.Errorf("Expected hash %s, got %s", doc.Hash, deserializedDoc.Hash)
	}
	if len(deserializedDoc.Metadata) != len(doc.Metadata) {
		t.Errorf("Expected metadata length %d, got %d", len(doc.Metadata), len(deserializedDoc.Metadata))
	}
	for key, value := range doc.Metadata {
		if deserializedDoc.Metadata[key] != value {
			t.Errorf("Expected metadata[%s] = %s, got %s", key, value, deserializedDoc.Metadata[key])
		}
	}
	if !deserializedDoc.CreatedAt.Equal(doc.CreatedAt) {
		t.Errorf("Expected created at %v, got %v", doc.CreatedAt, deserializedDoc.CreatedAt)
	}
	if !deserializedDoc.ProcessedAt.Equal(doc.ProcessedAt) {
		t.Errorf("Expected processed at %v, got %v", doc.ProcessedAt, deserializedDoc.ProcessedAt)
	}
}

func TestDocumentChunkSerialization(t *testing.T) {
	// Create test chunk
	chunk := &DocumentChunk{
		ID:         "chunk_123",
		DocumentID: "doc_456",
		ChunkIndex: 0,
		Content:    "This is a chunk of content",
		StartPos:   0,
		EndPos:     26,
		Embedding:  []float64{0.1, 0.2, 0.3, 0.4, 0.5},
		CreatedAt:  time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
	}

	// Test JSON serialization
	jsonData, err := json.Marshal(chunk)
	if err != nil {
		t.Fatalf("Failed to marshal chunk: %v", err)
	}

	// Test JSON deserialization
	var deserializedChunk DocumentChunk
	err = json.Unmarshal(jsonData, &deserializedChunk)
	if err != nil {
		t.Fatalf("Failed to unmarshal chunk: %v", err)
	}

	// Verify fields
	if deserializedChunk.ID != chunk.ID {
		t.Errorf("Expected ID %s, got %s", chunk.ID, deserializedChunk.ID)
	}
	if deserializedChunk.DocumentID != chunk.DocumentID {
		t.Errorf("Expected document ID %s, got %s", chunk.DocumentID, deserializedChunk.DocumentID)
	}
	if deserializedChunk.ChunkIndex != chunk.ChunkIndex {
		t.Errorf("Expected chunk index %d, got %d", chunk.ChunkIndex, deserializedChunk.ChunkIndex)
	}
	if deserializedChunk.Content != chunk.Content {
		t.Errorf("Expected content %s, got %s", chunk.Content, deserializedChunk.Content)
	}
	if deserializedChunk.StartPos != chunk.StartPos {
		t.Errorf("Expected start pos %d, got %d", chunk.StartPos, deserializedChunk.StartPos)
	}
	if deserializedChunk.EndPos != chunk.EndPos {
		t.Errorf("Expected end pos %d, got %d", chunk.EndPos, deserializedChunk.EndPos)
	}
	if len(deserializedChunk.Embedding) != len(chunk.Embedding) {
		t.Errorf("Expected embedding length %d, got %d", len(chunk.Embedding), len(deserializedChunk.Embedding))
	}
	for i, val := range chunk.Embedding {
		if deserializedChunk.Embedding[i] != val {
			t.Errorf("Expected embedding[%d] = %f, got %f", i, val, deserializedChunk.Embedding[i])
		}
	}
	if !deserializedChunk.CreatedAt.Equal(chunk.CreatedAt) {
		t.Errorf("Expected created at %v, got %v", chunk.CreatedAt, deserializedChunk.CreatedAt)
	}
}

func TestSearchResultSerialization(t *testing.T) {
	// Create test search result
	doc := &Document{
		ID:    "doc_123",
		Title: "Test Document",
	}
	chunk := &DocumentChunk{
		ID:         "chunk_123",
		DocumentID: "doc_123",
		Content:    "Test chunk content",
	}
	result := &SearchResult{
		Chunk:      chunk,
		Document:   doc,
		Similarity: 0.95,
	}

	// Test JSON serialization
	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal search result: %v", err)
	}

	// Test JSON deserialization
	var deserializedResult SearchResult
	err = json.Unmarshal(jsonData, &deserializedResult)
	if err != nil {
		t.Fatalf("Failed to unmarshal search result: %v", err)
	}

	// Verify fields
	if deserializedResult.Similarity != result.Similarity {
		t.Errorf("Expected similarity %f, got %f", result.Similarity, deserializedResult.Similarity)
	}
	if deserializedResult.Chunk == nil {
		t.Fatal("Chunk should not be nil")
	}
	if deserializedResult.Chunk.ID != chunk.ID {
		t.Errorf("Expected chunk ID %s, got %s", chunk.ID, deserializedResult.Chunk.ID)
	}
	if deserializedResult.Document == nil {
		t.Fatal("Document should not be nil")
	}
	if deserializedResult.Document.ID != doc.ID {
		t.Errorf("Expected document ID %s, got %s", doc.ID, deserializedResult.Document.ID)
	}
}

func TestRAGResponseSerialization(t *testing.T) {
	// Create test RAG response
	sources := []*SearchResult{
		{
			Chunk: &DocumentChunk{
				ID:      "chunk_1",
				Content: "First chunk",
			},
			Similarity: 0.9,
		},
		{
			Chunk: &DocumentChunk{
				ID:      "chunk_2",
				Content: "Second chunk",
			},
			Similarity: 0.8,
		},
	}

	response := &RAGResponse{
		Query:       "What is machine learning?",
		Answer:      "Machine learning is a subset of artificial intelligence.",
		Sources:     sources,
		ProcessTime: 2500 * time.Millisecond,
		CreatedAt:   time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
	}

	// Test JSON serialization
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal RAG response: %v", err)
	}

	// Test JSON deserialization
	var deserializedResponse RAGResponse
	err = json.Unmarshal(jsonData, &deserializedResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal RAG response: %v", err)
	}

	// Verify fields
	if deserializedResponse.Query != response.Query {
		t.Errorf("Expected query %s, got %s", response.Query, deserializedResponse.Query)
	}
	if deserializedResponse.Answer != response.Answer {
		t.Errorf("Expected answer %s, got %s", response.Answer, deserializedResponse.Answer)
	}
	if len(deserializedResponse.Sources) != len(response.Sources) {
		t.Errorf("Expected sources length %d, got %d", len(response.Sources), len(deserializedResponse.Sources))
	}
	for i, source := range response.Sources {
		if deserializedResponse.Sources[i].Similarity != source.Similarity {
			t.Errorf("Expected source[%d] similarity %f, got %f", i, source.Similarity, deserializedResponse.Sources[i].Similarity)
		}
	}
	if deserializedResponse.ProcessTime != response.ProcessTime {
		t.Errorf("Expected process time %v, got %v", response.ProcessTime, deserializedResponse.ProcessTime)
	}
	if !deserializedResponse.CreatedAt.Equal(response.CreatedAt) {
		t.Errorf("Expected created at %v, got %v", response.CreatedAt, deserializedResponse.CreatedAt)
	}
}

func TestEmbeddingRequestSerialization(t *testing.T) {
	// Create test embedding request
	req := &EmbeddingRequest{
		Texts: []string{"First text", "Second text", "Third text"},
		Model: "all-MiniLM-L6-v2",
	}

	// Test JSON serialization
	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal embedding request: %v", err)
	}

	// Test JSON deserialization
	var deserializedReq EmbeddingRequest
	err = json.Unmarshal(jsonData, &deserializedReq)
	if err != nil {
		t.Fatalf("Failed to unmarshal embedding request: %v", err)
	}

	// Verify fields
	if len(deserializedReq.Texts) != len(req.Texts) {
		t.Errorf("Expected texts length %d, got %d", len(req.Texts), len(deserializedReq.Texts))
	}
	for i, text := range req.Texts {
		if deserializedReq.Texts[i] != text {
			t.Errorf("Expected text[%d] = %s, got %s", i, text, deserializedReq.Texts[i])
		}
	}
	if deserializedReq.Model != req.Model {
		t.Errorf("Expected model %s, got %s", req.Model, deserializedReq.Model)
	}
}

func TestEmbeddingResponseSerialization(t *testing.T) {
	// Create test embedding response
	resp := &EmbeddingResponse{
		Embeddings: [][]float64{
			{0.1, 0.2, 0.3},
			{0.4, 0.5, 0.6},
			{0.7, 0.8, 0.9},
		},
		Model: "all-MiniLM-L6-v2",
	}

	// Test JSON serialization
	jsonData, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal embedding response: %v", err)
	}

	// Test JSON deserialization
	var deserializedResp EmbeddingResponse
	err = json.Unmarshal(jsonData, &deserializedResp)
	if err != nil {
		t.Fatalf("Failed to unmarshal embedding response: %v", err)
	}

	// Verify fields
	if len(deserializedResp.Embeddings) != len(resp.Embeddings) {
		t.Errorf("Expected embeddings length %d, got %d", len(resp.Embeddings), len(deserializedResp.Embeddings))
	}
	for i, embedding := range resp.Embeddings {
		if len(deserializedResp.Embeddings[i]) != len(embedding) {
			t.Errorf("Expected embedding[%d] length %d, got %d", i, len(embedding), len(deserializedResp.Embeddings[i]))
		}
		for j, val := range embedding {
			if deserializedResp.Embeddings[i][j] != val {
				t.Errorf("Expected embedding[%d][%d] = %f, got %f", i, j, val, deserializedResp.Embeddings[i][j])
			}
		}
	}
	if deserializedResp.Model != resp.Model {
		t.Errorf("Expected model %s, got %s", resp.Model, deserializedResp.Model)
	}
}

func TestLLMRequestSerialization(t *testing.T) {
	// Create test LLM request
	req := &LLMRequest{
		Prompt:      "What is the capital of France?",
		Model:       "llama2",
		Temperature: 0.7,
		MaxTokens:   512,
	}

	// Test JSON serialization
	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal LLM request: %v", err)
	}

	// Test JSON deserialization
	var deserializedReq LLMRequest
	err = json.Unmarshal(jsonData, &deserializedReq)
	if err != nil {
		t.Fatalf("Failed to unmarshal LLM request: %v", err)
	}

	// Verify fields
	if deserializedReq.Prompt != req.Prompt {
		t.Errorf("Expected prompt %s, got %s", req.Prompt, deserializedReq.Prompt)
	}
	if deserializedReq.Model != req.Model {
		t.Errorf("Expected model %s, got %s", req.Model, deserializedReq.Model)
	}
	if deserializedReq.Temperature != req.Temperature {
		t.Errorf("Expected temperature %f, got %f", req.Temperature, deserializedReq.Temperature)
	}
	if deserializedReq.MaxTokens != req.MaxTokens {
		t.Errorf("Expected max tokens %d, got %d", req.MaxTokens, deserializedReq.MaxTokens)
	}
}

func TestLLMResponseSerialization(t *testing.T) {
	// Create test LLM response
	resp := &LLMResponse{
		Response: "The capital of France is Paris.",
		Model:    "llama2",
	}

	// Test JSON serialization
	jsonData, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal LLM response: %v", err)
	}

	// Test JSON deserialization
	var deserializedResp LLMResponse
	err = json.Unmarshal(jsonData, &deserializedResp)
	if err != nil {
		t.Fatalf("Failed to unmarshal LLM response: %v", err)
	}

	// Verify fields
	if deserializedResp.Response != resp.Response {
		t.Errorf("Expected response %s, got %s", resp.Response, deserializedResp.Response)
	}
	if deserializedResp.Model != resp.Model {
		t.Errorf("Expected model %s, got %s", resp.Model, deserializedResp.Model)
	}
}

func TestEmptyAndNilFields(t *testing.T) {
	t.Run("DocumentWithEmptyFields", func(t *testing.T) {
		doc := &Document{
			ID:       "test",
			Metadata: map[string]string{},
		}

		jsonData, err := json.Marshal(doc)
		if err != nil {
			t.Fatalf("Failed to marshal document with empty fields: %v", err)
		}

		var deserializedDoc Document
		err = json.Unmarshal(jsonData, &deserializedDoc)
		if err != nil {
			t.Fatalf("Failed to unmarshal document with empty fields: %v", err)
		}

		if deserializedDoc.ID != doc.ID {
			t.Errorf("Expected ID %s, got %s", doc.ID, deserializedDoc.ID)
		}
		if deserializedDoc.Metadata == nil {
			t.Error("Metadata should not be nil after deserialization")
		}
	})

	t.Run("ChunkWithNilEmbedding", func(t *testing.T) {
		chunk := &DocumentChunk{
			ID:        "test",
			Embedding: nil,
		}

		jsonData, err := json.Marshal(chunk)
		if err != nil {
			t.Fatalf("Failed to marshal chunk with nil embedding: %v", err)
		}

		var deserializedChunk DocumentChunk
		err = json.Unmarshal(jsonData, &deserializedChunk)
		if err != nil {
			t.Fatalf("Failed to unmarshal chunk with nil embedding: %v", err)
		}

		if deserializedChunk.ID != chunk.ID {
			t.Errorf("Expected ID %s, got %s", chunk.ID, deserializedChunk.ID)
		}
		if deserializedChunk.Embedding != nil {
			t.Error("Embedding should be nil after deserialization")
		}
	})

	t.Run("EmbeddingRequestWithoutModel", func(t *testing.T) {
		req := &EmbeddingRequest{
			Texts: []string{"test"},
			// Model omitted
		}

		jsonData, err := json.Marshal(req)
		if err != nil {
			t.Fatalf("Failed to marshal embedding request without model: %v", err)
		}

		var deserializedReq EmbeddingRequest
		err = json.Unmarshal(jsonData, &deserializedReq)
		if err != nil {
			t.Fatalf("Failed to unmarshal embedding request without model: %v", err)
		}

		if len(deserializedReq.Texts) != 1 {
			t.Errorf("Expected 1 text, got %d", len(deserializedReq.Texts))
		}
		if deserializedReq.Model != "" {
			t.Errorf("Expected empty model, got %s", deserializedReq.Model)
		}
	})
}