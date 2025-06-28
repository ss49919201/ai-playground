package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"yesql-account-system/internal/account"
	"yesql-account-system/internal/auth"

	"github.com/gorilla/mux"
)

type Handler struct {
	accountService *account.Service
	authService    *auth.Service
}

func NewHandler(accountService *account.Service, authService *auth.Service) *Handler {
	return &Handler{
		accountService: accountService,
		authService:    authService,
	}
}

func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req account.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	acc, err := h.accountService.CreateAccount(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create account: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(acc)
}

func (h *Handler) GetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["id"]

	acc, err := h.accountService.GetAccount(accountID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get account: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(acc)
}

func (h *Handler) ListAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.accountService.ListAccounts()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list accounts: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

func (h *Handler) Deposit(w http.ResponseWriter, r *http.Request) {
	var req account.DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	transaction, err := h.accountService.Deposit(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to deposit: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	var req account.WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	transaction, err := h.accountService.Withdraw(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to withdraw: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req account.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	transaction, err := h.accountService.Transfer(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to transfer: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

func (h *Handler) GetAccountTransactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["id"]

	transactions, err := h.accountService.GetAccountTransactions(accountID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get transactions: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

// Auth handlers
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req auth.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to register user: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	loginResponse, err := h.authService.Login(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to login: %v", err), http.StatusUnauthorized)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    loginResponse.SessionID,
		Expires:  loginResponse.ExpiresAt,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionID := auth.GetSessionIDFromContext(r.Context())
	if sessionID == "" {
		http.Error(w, "No active session", http.StatusBadRequest)
		return
	}

	err := h.authService.Logout(sessionID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to logout: %v", err), http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout successful"))
}

func (h *Handler) GetUserAccounts(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	accounts, err := h.authService.GetUserAccounts(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user accounts: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}