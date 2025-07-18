package handler

import (
	"encoding/json"
	"go-memo-app/internal/usecase"
	"net/http"
	"strconv"
	"strings"
)

func NewMemoHandler(
	getAllMemos usecase.GetAllMemosFunc,
	getMemoByID usecase.GetMemoByIDFunc,
	createMemo usecase.CreateMemoFunc,
	updateMemo usecase.UpdateMemoFunc,
	deleteMemo usecase.DeleteMemoFunc,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGet(w, r, getAllMemos, getMemoByID)
		case http.MethodPost:
			handlePost(w, r, createMemo)
		case http.MethodPut:
			handlePut(w, r, updateMemo)
		case http.MethodDelete:
			handleDelete(w, r, deleteMemo)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func handleGet(w http.ResponseWriter, r *http.Request, getAllMemos usecase.GetAllMemosFunc, getMemoByID usecase.GetMemoByIDFunc) {
	path := strings.TrimPrefix(r.URL.Path, "/api/memos")
	
	if path == "" || path == "/" {
		memos := getAllMemos()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(memos)
		return
	}

	idStr := strings.TrimPrefix(path, "/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid memo ID", http.StatusBadRequest)
		return
	}

	memo, err := getMemoByID(id)
	if err != nil {
		http.Error(w, "Memo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(memo)
}

func handlePost(w http.ResponseWriter, r *http.Request, createMemo usecase.CreateMemoFunc) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	memo, err := createMemo(req.Title, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(memo)
}

func handlePut(w http.ResponseWriter, r *http.Request, updateMemo usecase.UpdateMemoFunc) {
	path := strings.TrimPrefix(r.URL.Path, "/api/memos/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid memo ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	memo, err := updateMemo(id, req.Title, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(memo)
}

func handleDelete(w http.ResponseWriter, r *http.Request, deleteMemo usecase.DeleteMemoFunc) {
	path := strings.TrimPrefix(r.URL.Path, "/api/memos/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid memo ID", http.StatusBadRequest)
		return
	}

	if err := deleteMemo(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}