package dao

import (
	"testing"
	"time"

	"go-memo-app/pkg/model"
)

func TestNewInMemoryGetAllMemos(t *testing.T) {
	memos = make(map[int]*model.Memo)
	nextID = 1

	getAllMemos := NewInMemoryGetAllMemos()
	
	result := getAllMemos()
	if len(result) != 0 {
		t.Errorf("Expected empty slice, got %d memos", len(result))
	}

	memos[1] = &model.Memo{
		ID:      1,
		Title:   "Test Memo",
		Content: "Test Content",
	}

	result = getAllMemos()
	if len(result) != 1 {
		t.Errorf("Expected 1 memo, got %d memos", len(result))
	}
	if result[0].Title != "Test Memo" {
		t.Errorf("Expected title 'Test Memo', got %s", result[0].Title)
	}
}

func TestNewInMemoryGetMemoByID(t *testing.T) {
	memos = make(map[int]*model.Memo)
	nextID = 1

	getMemoByID := NewInMemoryGetMemoByID()

	_, err := getMemoByID(1)
	if err != ErrMemoNotFound {
		t.Errorf("Expected ErrMemoNotFound, got %v", err)
	}

	memos[1] = &model.Memo{
		ID:      1,
		Title:   "Test Memo",
		Content: "Test Content",
	}

	memo, err := getMemoByID(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if memo.Title != "Test Memo" {
		t.Errorf("Expected title 'Test Memo', got %s", memo.Title)
	}
}

func TestNewInMemoryCreateMemo(t *testing.T) {
	memos = make(map[int]*model.Memo)
	nextID = 1

	createMemo := NewInMemoryCreateMemo()

	memo, err := createMemo("Test Title", "Test Content")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if memo.ID != 1 {
		t.Errorf("Expected ID 1, got %d", memo.ID)
	}
	if memo.Title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got %s", memo.Title)
	}
	if memo.Content != "Test Content" {
		t.Errorf("Expected content 'Test Content', got %s", memo.Content)
	}
	if memo.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
	if memo.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}

	if len(memos) != 1 {
		t.Errorf("Expected 1 memo in storage, got %d", len(memos))
	}
	if nextID != 2 {
		t.Errorf("Expected nextID to be 2, got %d", nextID)
	}
}

func TestNewInMemoryUpdateMemo(t *testing.T) {
	memos = make(map[int]*model.Memo)
	nextID = 1

	updateMemo := NewInMemoryUpdateMemo()

	_, err := updateMemo(1, "Updated Title", "Updated Content")
	if err != ErrMemoNotFound {
		t.Errorf("Expected ErrMemoNotFound, got %v", err)
	}

	createdAt := time.Now().Add(-time.Hour)
	memos[1] = &model.Memo{
		ID:        1,
		Title:     "Original Title",
		Content:   "Original Content",
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}

	memo, err := updateMemo(1, "Updated Title", "Updated Content")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if memo.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got %s", memo.Title)
	}
	if memo.Content != "Updated Content" {
		t.Errorf("Expected content 'Updated Content', got %s", memo.Content)
	}
	if memo.CreatedAt != createdAt {
		t.Error("Expected CreatedAt to remain unchanged")
	}
	if !memo.UpdatedAt.After(createdAt) {
		t.Error("Expected UpdatedAt to be updated")
	}
}

func TestNewInMemoryDeleteMemo(t *testing.T) {
	memos = make(map[int]*model.Memo)
	nextID = 1

	deleteMemo := NewInMemoryDeleteMemo()

	err := deleteMemo(1)
	if err != ErrMemoNotFound {
		t.Errorf("Expected ErrMemoNotFound, got %v", err)
	}

	memos[1] = &model.Memo{
		ID:      1,
		Title:   "Test Memo",
		Content: "Test Content",
	}

	err = deleteMemo(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(memos) != 0 {
		t.Errorf("Expected 0 memos in storage, got %d", len(memos))
	}
}