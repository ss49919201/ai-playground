package usecase

import (
	"errors"
	"testing"
	"time"

	"go-memo-app/pkg/model"
)

func TestNewGetAllMemos(t *testing.T) {
	mockGetAllMemos := func() []model.Memo {
		return []model.Memo{
			{ID: 1, Title: "Memo 1", Content: "Content 1"},
			{ID: 2, Title: "Memo 2", Content: "Content 2"},
		}
	}

	getAllMemos := NewGetAllMemos(mockGetAllMemos)
	result := getAllMemos()

	if len(result) != 2 {
		t.Errorf("Expected 2 memos, got %d", len(result))
	}
	if result[0].Title != "Memo 1" {
		t.Errorf("Expected title 'Memo 1', got %s", result[0].Title)
	}
}

func TestNewGetMemoByID(t *testing.T) {
	mockGetMemoByID := func(id int) (*model.Memo, error) {
		if id == 1 {
			return &model.Memo{ID: 1, Title: "Test Memo", Content: "Test Content"}, nil
		}
		return nil, errors.New("memo not found")
	}

	getMemoByID := NewGetMemoByID(mockGetMemoByID)

	memo, err := getMemoByID(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if memo.Title != "Test Memo" {
		t.Errorf("Expected title 'Test Memo', got %s", memo.Title)
	}

	_, err = getMemoByID(999)
	if err == nil {
		t.Error("Expected error for non-existent memo")
	}
}

func TestNewCreateMemo(t *testing.T) {
	mockCreateMemo := func(title, content string) (*model.Memo, error) {
		return &model.Memo{
			ID:        1,
			Title:     title,
			Content:   content,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil
	}

	createMemo := NewCreateMemo(mockCreateMemo)

	memo, err := createMemo("Test Title", "Test Content")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if memo.Title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got %s", memo.Title)
	}

	_, err = createMemo("", "Test Content")
	if err != ErrInvalidInput {
		t.Errorf("Expected ErrInvalidInput, got %v", err)
	}
}

func TestNewUpdateMemo(t *testing.T) {
	mockUpdateMemo := func(id int, title, content string) (*model.Memo, error) {
		if id == 1 {
			return &model.Memo{
				ID:        1,
				Title:     title,
				Content:   content,
				CreatedAt: time.Now().Add(-time.Hour),
				UpdatedAt: time.Now(),
			}, nil
		}
		return nil, errors.New("memo not found")
	}

	updateMemo := NewUpdateMemo(mockUpdateMemo)

	memo, err := updateMemo(1, "Updated Title", "Updated Content")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if memo.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got %s", memo.Title)
	}

	_, err = updateMemo(1, "", "Updated Content")
	if err != ErrInvalidInput {
		t.Errorf("Expected ErrInvalidInput, got %v", err)
	}

	_, err = updateMemo(999, "Updated Title", "Updated Content")
	if err == nil {
		t.Error("Expected error for non-existent memo")
	}
}

func TestNewDeleteMemo(t *testing.T) {
	mockDeleteMemo := func(id int) error {
		if id == 1 {
			return nil
		}
		return errors.New("memo not found")
	}

	deleteMemo := NewDeleteMemo(mockDeleteMemo)

	err := deleteMemo(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = deleteMemo(999)
	if err == nil {
		t.Error("Expected error for non-existent memo")
	}
}