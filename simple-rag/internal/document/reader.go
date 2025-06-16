package document

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"simple-rag/pkg/types"
)

// Reader handles document reading and processing
type Reader struct {
	chunkSize    int
	chunkOverlap int
}

// NewReader creates a new document reader
func NewReader(chunkSize, chunkOverlap int) *Reader {
	return &Reader{
		chunkSize:    chunkSize,
		chunkOverlap: chunkOverlap,
	}
}

// ReadDocument reads a document from file path and returns a Document
func (r *Reader) ReadDocument(filePath string) (*types.Document, error) {
	// Get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Calculate file hash
	hash := fmt.Sprintf("%x", md5.Sum(content))

	// Extract file info
	fileName := filepath.Base(filePath)
	fileExt := strings.ToLower(filepath.Ext(filePath))
	if fileExt != "" {
		fileExt = fileExt[1:] // Remove the dot
	}

	// Create document
	doc := &types.Document{
		ID:          generateDocumentID(filePath, hash),
		Title:       fileName,
		Content:     string(content),
		FilePath:    filePath,
		FileType:    fileExt,
		FileSize:    fileInfo.Size(),
		Hash:        hash,
		Metadata:    make(map[string]string),
		CreatedAt:   fileInfo.ModTime(),
		ProcessedAt: time.Now(),
	}

	// Add basic metadata
	doc.Metadata["original_name"] = fileName
	doc.Metadata["file_extension"] = fileExt
	doc.Metadata["content_length"] = fmt.Sprintf("%d", len(doc.Content))

	return doc, nil
}

// ChunkDocument splits a document into chunks
func (r *Reader) ChunkDocument(doc *types.Document) ([]*types.DocumentChunk, error) {
	content := doc.Content
	if len(content) == 0 {
		return nil, fmt.Errorf("document content is empty")
	}

	var chunks []*types.DocumentChunk
	chunkIndex := 0
	startPos := 0

	for startPos < len(content) {
		endPos := startPos + r.chunkSize
		if endPos > len(content) {
			endPos = len(content)
		}

		// Try to break at word boundary
		if endPos < len(content) {
			for i := endPos; i > startPos; i-- {
				if content[i] == ' ' || content[i] == '\n' || content[i] == '\t' {
					endPos = i
					break
				}
			}
		}

		chunkContent := content[startPos:endPos]
		chunkContent = strings.TrimSpace(chunkContent)

		if len(chunkContent) > 0 {
			chunk := &types.DocumentChunk{
				ID:         generateChunkID(doc.ID, chunkIndex),
				DocumentID: doc.ID,
				ChunkIndex: chunkIndex,
				Content:    chunkContent,
				StartPos:   startPos,
				EndPos:     endPos,
				Embedding:  nil, // Will be set later
				CreatedAt:  time.Now(),
			}
			chunks = append(chunks, chunk)
			chunkIndex++
		}

		// Move to next chunk with overlap
		nextStart := endPos - r.chunkOverlap
		if nextStart <= startPos {
			nextStart = startPos + 1
		}
		startPos = nextStart
	}

	return chunks, nil
}

// generateDocumentID generates a unique ID for a document
func generateDocumentID(filePath, hash string) string {
	return fmt.Sprintf("doc_%x", md5.Sum([]byte(filePath+hash)))[:16]
}

// generateChunkID generates a unique ID for a chunk
func generateChunkID(documentID string, chunkIndex int) string {
	return fmt.Sprintf("%s_chunk_%d", documentID, chunkIndex)
}

// IsSupported checks if the file type is supported
func IsSupported(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	supportedExts := map[string]bool{
		".txt":  true,
		".md":   true,
		".text": true,
	}
	return supportedExts[ext]
}
