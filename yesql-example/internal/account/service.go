package account

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"yesql-account-system/internal/yesql"
)

type Service struct {
	db     *sql.DB
	loader *yesql.QueryLoader
}

func NewService(db *sql.DB, loader *yesql.QueryLoader) *Service {
	return &Service{
		db:     db,
		loader: loader,
	}
}

func (s *Service) CreateAccount(req CreateAccountRequest) (*Account, error) {
	query, err := s.loader.GetQuery("create_account")
	if err != nil {
		return nil, fmt.Errorf("failed to get create_account query: %w", err)
	}

	_, err = s.db.Exec(query, req.AccountID, req.AccountName, req.InitialDeposit)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return s.GetAccount(req.AccountID)
}

func (s *Service) GetAccount(accountID string) (*Account, error) {
	query, err := s.loader.GetQuery("get_account_by_id")
	if err != nil {
		return nil, fmt.Errorf("failed to get get_account_by_id query: %w", err)
	}

	var account Account
	err = s.db.QueryRow(query, accountID).Scan(
		&account.AccountID,
		&account.AccountName,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &account, nil
}

func (s *Service) ListAccounts() ([]Account, error) {
	query, err := s.loader.GetQuery("list_accounts")
	if err != nil {
		return nil, fmt.Errorf("failed to get list_accounts query: %w", err)
	}

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		var account Account
		err := rows.Scan(
			&account.AccountID,
			&account.AccountName,
			&account.Balance,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *Service) Deposit(req DepositRequest) (*Transaction, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	updateQuery, err := s.loader.GetQuery("deposit_update_balance")
	if err != nil {
		return nil, fmt.Errorf("failed to get deposit_update_balance query: %w", err)
	}

	result, err := tx.Exec(updateQuery, req.Amount, req.AccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("account not found")
	}

	createTxQuery, err := s.loader.GetQuery("deposit_create_transaction")
	if err != nil {
		return nil, fmt.Errorf("failed to get deposit_create_transaction query: %w", err)
	}

	transactionID := uuid.New().String()
	var transaction Transaction
	err = tx.QueryRow(createTxQuery, transactionID, req.AccountID, req.Amount, req.Description).Scan(
		&transaction.TransactionID,
		&transaction.FromAccount,
		&transaction.ToAccount,
		&transaction.TransactionType,
		&transaction.Amount,
		&transaction.Description,
		&transaction.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &transaction, nil
}

func (s *Service) Withdraw(req WithdrawRequest) (*Transaction, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	updateQuery, err := s.loader.GetQuery("withdraw_update_balance")
	if err != nil {
		return nil, fmt.Errorf("failed to get withdraw_update_balance query: %w", err)
	}

	result, err := tx.Exec(updateQuery, req.Amount, req.AccountID, req.Amount)
	if err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("insufficient funds or account not found")
	}

	createTxQuery, err := s.loader.GetQuery("withdraw_create_transaction")
	if err != nil {
		return nil, fmt.Errorf("failed to get withdraw_create_transaction query: %w", err)
	}

	transactionID := uuid.New().String()
	var transaction Transaction
	err = tx.QueryRow(createTxQuery, transactionID, req.AccountID, req.Amount, req.Description).Scan(
		&transaction.TransactionID,
		&transaction.FromAccount,
		&transaction.ToAccount,
		&transaction.TransactionType,
		&transaction.Amount,
		&transaction.Description,
		&transaction.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &transaction, nil
}

func (s *Service) Transfer(req TransferRequest) (*Transaction, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	balanceQuery, err := s.loader.GetQuery("get_account_balance")
	if err != nil {
		return nil, fmt.Errorf("failed to get get_account_balance query: %w", err)
	}

	var fromBalance, toBalance float64
	err = tx.QueryRow(balanceQuery, req.FromAccountID).Scan(&fromBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to get from account balance: %w", err)
	}

	err = tx.QueryRow(balanceQuery, req.ToAccountID).Scan(&toBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to get to account balance: %w", err)
	}

	if fromBalance < req.Amount {
		return nil, fmt.Errorf("insufficient funds: current balance %.2f, requested %.2f", fromBalance, req.Amount)
	}

	updateQuery, err := s.loader.GetQuery("update_account_balance")
	if err != nil {
		return nil, fmt.Errorf("failed to get update_account_balance query: %w", err)
	}

	_, err = tx.Exec(updateQuery, fromBalance-req.Amount, req.FromAccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to update from account balance: %w", err)
	}

	_, err = tx.Exec(updateQuery, toBalance+req.Amount, req.ToAccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to update to account balance: %w", err)
	}

	transactionID := uuid.New().String()
	createTxQuery, err := s.loader.GetQuery("create_transaction")
	if err != nil {
		return nil, fmt.Errorf("failed to get create_transaction query: %w", err)
	}

	_, err = tx.Exec(createTxQuery, transactionID, req.FromAccountID, req.ToAccountID, "transfer", req.Amount, req.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return s.GetTransaction(transactionID)
}

func (s *Service) GetTransaction(transactionID string) (*Transaction, error) {
	query, err := s.loader.GetQuery("get_transaction_by_id")
	if err != nil {
		return nil, fmt.Errorf("failed to get get_transaction_by_id query: %w", err)
	}

	var transaction Transaction
	err = s.db.QueryRow(query, transactionID).Scan(
		&transaction.TransactionID,
		&transaction.FromAccount,
		&transaction.ToAccount,
		&transaction.TransactionType,
		&transaction.Amount,
		&transaction.Description,
		&transaction.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &transaction, nil
}

func (s *Service) GetAccountTransactions(accountID string) ([]Transaction, error) {
	query, err := s.loader.GetQuery("get_account_transactions")
	if err != nil {
		return nil, fmt.Errorf("failed to get get_account_transactions query: %w", err)
	}

	rows, err := s.db.Query(query, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account transactions: %w", err)
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var transaction Transaction
		err := rows.Scan(
			&transaction.TransactionID,
			&transaction.FromAccount,
			&transaction.ToAccount,
			&transaction.TransactionType,
			&transaction.Amount,
			&transaction.Description,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}