package account

import (
	"time"
)

type Account struct {
	AccountID   string    `json:"account_id"`
	AccountName string    `json:"account_name"`
	Balance     float64   `json:"balance"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Transaction struct {
	TransactionID   string    `json:"transaction_id"`
	FromAccount     *string   `json:"from_account"`
	ToAccount       *string   `json:"to_account"`
	TransactionType string    `json:"transaction_type"`
	Amount          float64   `json:"amount"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
}

type TransferRequest struct {
	FromAccountID string  `json:"from_account_id"`
	ToAccountID   string  `json:"to_account_id"`
	Amount        float64 `json:"amount"`
	Description   string  `json:"description"`
}

type DepositRequest struct {
	AccountID   string  `json:"account_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

type WithdrawRequest struct {
	AccountID   string  `json:"account_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

type CreateAccountRequest struct {
	AccountID    string  `json:"account_id"`
	AccountName  string  `json:"account_name"`
	InitialDeposit float64 `json:"initial_deposit"`
}