package yesql

import (
	"os"
	"path/filepath"
	"testing"
)

func TestQueryLoader(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "yesql-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testSQL := `-- name: test_query_1
SELECT * FROM accounts WHERE id = $1;

-- name: test_query_2
INSERT INTO accounts (id, name) VALUES ($1, $2);

-- name: test_query_3
UPDATE accounts SET balance = $2 WHERE id = $1;`

	testFile := filepath.Join(tempDir, "test.sql")
	err = os.WriteFile(testFile, []byte(testSQL), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	loader := NewQueryLoader()
	err = loader.LoadQueriesFromDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to load queries: %v", err)
	}

	query1, err := loader.GetQuery("test_query_1")
	if err != nil {
		t.Fatalf("Failed to get test_query_1: %v", err)
	}

	expectedQuery1 := "SELECT * FROM accounts WHERE id = $1;"
	if query1 != expectedQuery1 {
		t.Errorf("Expected query '%s', got '%s'", expectedQuery1, query1)
	}

	query2, err := loader.GetQuery("test_query_2")
	if err != nil {
		t.Fatalf("Failed to get test_query_2: %v", err)
	}

	expectedQuery2 := "INSERT INTO accounts (id, name) VALUES ($1, $2);"
	if query2 != expectedQuery2 {
		t.Errorf("Expected query '%s', got '%s'", expectedQuery2, query2)
	}

	_, err = loader.GetQuery("nonexistent_query")
	if err == nil {
		t.Error("Expected error for nonexistent query, but got none")
	}

	queries := loader.ListQueries()
	if len(queries) != 3 {
		t.Errorf("Expected 3 queries, got %d", len(queries))
	}
}