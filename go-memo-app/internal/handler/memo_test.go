package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-memo-app/pkg/model"
)

func TestNewMemoHandler_GET_All(t *testing.T) {
	mockGetAllMemos := func() []model.Memo {
		return []model.Memo{
			{ID: 1, Title: "Memo 1", Content: "Content 1"},
			{ID: 2, Title: "Memo 2", Content: "Content 2"},
		}
	}

	mockGetMemoByID := func(id int) (*model.Memo, error) {
		return nil, errors.New("not used")
	}

	handler := NewMemoHandler(
		mockGetAllMemos,
		mockGetMemoByID,
		nil, nil, nil,
	)

	req := httptest.NewRequest(http.MethodGet, "/api/memos", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var memos []model.Memo
	err := json.NewDecoder(w.Body).Decode(&memos)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if len(memos) != 2 {
		t.Errorf("Expected 2 memos, got %d", len(memos))
	}
}

func TestNewMemoHandler_GET_ByID(t *testing.T) {
	mockGetAllMemos := func() []model.Memo {
		return nil
	}

	mockGetMemoByID := func(id int) (*model.Memo, error) {
		if id == 1 {
			return &model.Memo{ID: 1, Title: "Test Memo", Content: "Test Content"}, nil
		}
		return nil, errors.New("memo not found")
	}

	handler := NewMemoHandler(
		mockGetAllMemos,
		mockGetMemoByID,
		nil, nil, nil,
	)

	req := httptest.NewRequest(http.MethodGet, "/api/memos/1", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var memo model.Memo
	err := json.NewDecoder(w.Body).Decode(&memo)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if memo.Title != "Test Memo" {
		t.Errorf("Expected title 'Test Memo', got %s", memo.Title)
	}
}

func TestNewMemoHandler_GET_ByID_NotFound(t *testing.T) {
	mockGetAllMemos := func() []model.Memo {
		return nil
	}

	mockGetMemoByID := func(id int) (*model.Memo, error) {
		return nil, errors.New("memo not found")
	}

	handler := NewMemoHandler(
		mockGetAllMemos,
		mockGetMemoByID,
		nil, nil, nil,
	)

	req := httptest.NewRequest(http.MethodGet, "/api/memos/999", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestNewMemoHandler_POST(t *testing.T) {
	mockCreateMemo := func(title, content string) (*model.Memo, error) {
		return &model.Memo{
			ID:        1,
			Title:     title,
			Content:   content,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil
	}

	handler := NewMemoHandler(
		nil, nil,
		mockCreateMemo,
		nil, nil,
	)

	reqBody := map[string]string{
		"title":   "New Memo",
		"content": "New Content",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/memos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var memo model.Memo
	err := json.NewDecoder(w.Body).Decode(&memo)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if memo.Title != "New Memo" {
		t.Errorf("Expected title 'New Memo', got %s", memo.Title)
	}
}

func TestNewMemoHandler_PUT(t *testing.T) {
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

	handler := NewMemoHandler(
		nil, nil, nil,
		mockUpdateMemo,
		nil,
	)

	reqBody := map[string]string{
		"title":   "Updated Memo",
		"content": "Updated Content",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/memos/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var memo model.Memo
	err := json.NewDecoder(w.Body).Decode(&memo)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if memo.Title != "Updated Memo" {
		t.Errorf("Expected title 'Updated Memo', got %s", memo.Title)
	}
}

func TestNewMemoHandler_DELETE(t *testing.T) {
	mockDeleteMemo := func(id int) error {
		if id == 1 {
			return nil
		}
		return errors.New("memo not found")
	}

	handler := NewMemoHandler(
		nil, nil, nil, nil,
		mockDeleteMemo,
	)

	req := httptest.NewRequest(http.MethodDelete, "/api/memos/1", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestNewMemoHandler_DELETE_NotFound(t *testing.T) {
	mockDeleteMemo := func(id int) error {
		return errors.New("memo not found")
	}

	handler := NewMemoHandler(
		nil, nil, nil, nil,
		mockDeleteMemo,
	)

	req := httptest.NewRequest(http.MethodDelete, "/api/memos/999", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestNewMemoHandler_InvalidMethod(t *testing.T) {
	handler := NewMemoHandler(nil, nil, nil, nil, nil)

	req := httptest.NewRequest(http.MethodPatch, "/api/memos", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestNewMemoHandler_InvalidJSON(t *testing.T) {
	mockCreateMemo := func(title, content string) (*model.Memo, error) {
		return nil, errors.New("should not be called")
	}

	handler := NewMemoHandler(
		nil, nil,
		mockCreateMemo,
		nil, nil,
	)

	req := httptest.NewRequest(http.MethodPost, "/api/memos", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestNewMemoHandler_InvalidID(t *testing.T) {
	mockGetMemoByID := func(id int) (*model.Memo, error) {
		return nil, errors.New("should not be called")
	}

	handler := NewMemoHandler(
		nil,
		mockGetMemoByID,
		nil, nil, nil,
	)

	req := httptest.NewRequest(http.MethodGet, "/api/memos/invalid", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}