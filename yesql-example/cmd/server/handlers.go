package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yesql-account-system/internal/account"

	"github.com/gorilla/mux"
)

type Handler struct {
	accountService *account.Service
}

func NewHandler(accountService *account.Service) *Handler {
	return &Handler{
		accountService: accountService,
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