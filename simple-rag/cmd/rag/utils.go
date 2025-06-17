package main

import (
	"fmt"
	"time"

	"simple-rag/pkg/types"
)

// printResponse prints a RAG response in a formatted way
func printResponse(response *types.RAGResponse) {
	fmt.Println(repeatString("=", 60))
	fmt.Printf("Query: %s\n", response.Query)
	fmt.Println(repeatString("-", 60))
	fmt.Printf("Answer: %s\n", response.Answer)
	fmt.Println(repeatString("-", 60))
	
	if len(response.Sources) > 0 {
		fmt.Println("Sources:")
		for i, source := range response.Sources {
			fmt.Printf("\n%d. Document: %s (Similarity: %.3f)\n", 
				i+1, source.Document.Title, source.Similarity)
			fmt.Printf("   Content: %s...\n", truncateString(source.Chunk.Content, 100))
			fmt.Printf("   File: %s\n", source.Document.FilePath)
		}
	}
	
	fmt.Println(repeatString("-", 60))
	fmt.Printf("Process Time: %v\n", response.ProcessTime)
	fmt.Printf("Generated At: %s\n", response.CreatedAt.Format(time.RFC3339))
	fmt.Println(repeatString("=", 60))
}

// printDocuments prints a list of documents
func printDocuments(documents []*types.Document) {
	if len(documents) == 0 {
		fmt.Println("No documents found.")
		return
	}

	fmt.Println("Documents in the system:")
	fmt.Println(repeatString("-", 80))
	for i, doc := range documents {
		fmt.Printf("%d. Title: %s\n", i+1, doc.Title)
		fmt.Printf("   ID: %s\n", doc.ID)
		fmt.Printf("   File: %s\n", doc.FilePath)
		fmt.Printf("   Type: %s\n", doc.FileType)
		fmt.Printf("   Size: %d bytes\n", doc.FileSize)
		fmt.Printf("   Content: %s...\n", truncateString(doc.Content, 100))
		fmt.Printf("   Created: %s\n", doc.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("   Processed: %s\n", doc.ProcessedAt.Format("2006-01-02 15:04:05"))
		if i < len(documents)-1 {
			fmt.Println(repeatString("-", 80))
		}
	}
	fmt.Println(repeatString("-", 80))
}

// truncateString truncates a string to a maximum length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// repeatString repeats a string n times (Go doesn't have built-in string repeat)
func repeatString(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

// Use repeatString for the separator lines
func init() {
	// Replace the "*" and "-" repetition with function calls
	// This will be handled in the print functions
}