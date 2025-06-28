package account

import (
	"database/sql"
	"os"
	"testing"
	"yesql-account-system/internal/db"
	"yesql-account-system/internal/yesql"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) (*sql.DB, *Service) {
	database, err := db.NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	err = database.InitSchema("../../sql/schema.sql")
	if err != nil {
		t.Fatalf("Failed to initialize schema: %v", err)
	}

	loader := yesql.NewQueryLoader()
	err = loader.LoadQueriesFromDir("../../sql/queries")
	if err != nil {
		t.Fatalf("Failed to load queries: %v", err)
	}

	service := NewService(database.GetConn(), loader)
	return database.GetConn(), service
}

func TestCreateAccount(t *testing.T) {
	dbConn, service := setupTestDB(t)
	defer dbConn.Close()

	req := CreateAccountRequest{
		AccountID:      "test-account-1",
		AccountName:    "Test Account",
		InitialDeposit: 1000.0,
	}

	account, err := service.CreateAccount(req)
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}

	if account.AccountID != req.AccountID {
		t.Errorf("Expected account ID %s, got %s", req.AccountID, account.AccountID)
	}

	if account.AccountName != req.AccountName {
		t.Errorf("Expected account name %s, got %s", req.AccountName, account.AccountName)
	}

	if account.Balance != req.InitialDeposit {
		t.Errorf("Expected balance %f, got %f", req.InitialDeposit, account.Balance)
	}
}

func TestGetAccount(t *testing.T) {
	dbConn, service := setupTestDB(t)
	defer dbConn.Close()

	req := CreateAccountRequest{
		AccountID:      "test-account-2",
		AccountName:    "Test Account 2",
		InitialDeposit: 500.0,
	}

	_, err := service.CreateAccount(req)
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}

	account, err := service.GetAccount(req.AccountID)
	if err != nil {
		t.Fatalf("Failed to get account: %v", err)
	}

	if account.AccountID != req.AccountID {
		t.Errorf("Expected account ID %s, got %s", req.AccountID, account.AccountID)
	}
}

func TestDeposit(t *testing.T) {
	dbConn, service := setupTestDB(t)
	defer dbConn.Close()

	createReq := CreateAccountRequest{
		AccountID:      "test-account-3",
		AccountName:    "Test Account 3",
		InitialDeposit: 100.0,
	}

	_, err := service.CreateAccount(createReq)
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}

	depositReq := DepositRequest{
		AccountID:   "test-account-3",
		Amount:      200.0,
		Description: "Test deposit",
	}

	transaction, err := service.Deposit(depositReq)
	if err != nil {
		t.Fatalf("Failed to deposit: %v", err)
	}

	if transaction.TransactionType != "deposit" {
		t.Errorf("Expected transaction type 'deposit', got %s", transaction.TransactionType)
	}

	if transaction.Amount != depositReq.Amount {
		t.Errorf("Expected amount %f, got %f", depositReq.Amount, transaction.Amount)
	}

	account, err := service.GetAccount("test-account-3")
	if err != nil {
		t.Fatalf("Failed to get account: %v", err)
	}

	expectedBalance := 300.0
	if account.Balance != expectedBalance {
		t.Errorf("Expected balance %f, got %f", expectedBalance, account.Balance)
	}
}

func TestWithdraw(t *testing.T) {
	dbConn, service := setupTestDB(t)
	defer dbConn.Close()

	createReq := CreateAccountRequest{
		AccountID:      "test-account-4",
		AccountName:    "Test Account 4",
		InitialDeposit: 500.0,
	}

	_, err := service.CreateAccount(createReq)
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}

	withdrawReq := WithdrawRequest{
		AccountID:   "test-account-4",
		Amount:      100.0,
		Description: "Test withdrawal",
	}

	transaction, err := service.Withdraw(withdrawReq)
	if err != nil {
		t.Fatalf("Failed to withdraw: %v", err)
	}

	if transaction.TransactionType != "withdrawal" {
		t.Errorf("Expected transaction type 'withdrawal', got %s", transaction.TransactionType)
	}

	account, err := service.GetAccount("test-account-4")
	if err != nil {
		t.Fatalf("Failed to get account: %v", err)
	}

	expectedBalance := 400.0
	if account.Balance != expectedBalance {
		t.Errorf("Expected balance %f, got %f", expectedBalance, account.Balance)
	}
}

func TestWithdrawInsufficientFunds(t *testing.T) {
	dbConn, service := setupTestDB(t)
	defer dbConn.Close()

	createReq := CreateAccountRequest{
		AccountID:      "test-account-5",
		AccountName:    "Test Account 5",
		InitialDeposit: 100.0,
	}

	_, err := service.CreateAccount(createReq)
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}

	withdrawReq := WithdrawRequest{
		AccountID:   "test-account-5",
		Amount:      200.0,
		Description: "Test insufficient withdrawal",
	}

	_, err = service.Withdraw(withdrawReq)
	if err == nil {
		t.Error("Expected error for insufficient funds, but got none")
	}
}

func TestTransfer(t *testing.T) {
	dbConn, service := setupTestDB(t)
	defer dbConn.Close()

	fromReq := CreateAccountRequest{
		AccountID:      "test-from-account",
		AccountName:    "From Account",
		InitialDeposit: 1000.0,
	}

	toReq := CreateAccountRequest{
		AccountID:      "test-to-account",
		AccountName:    "To Account",
		InitialDeposit: 0.0,
	}

	_, err := service.CreateAccount(fromReq)
	if err != nil {
		t.Fatalf("Failed to create from account: %v", err)
	}

	_, err = service.CreateAccount(toReq)
	if err != nil {
		t.Fatalf("Failed to create to account: %v", err)
	}

	transferReq := TransferRequest{
		FromAccountID: "test-from-account",
		ToAccountID:   "test-to-account",
		Amount:        300.0,
		Description:   "Test transfer",
	}

	transaction, err := service.Transfer(transferReq)
	if err != nil {
		t.Fatalf("Failed to transfer: %v", err)
	}

	if transaction.TransactionType != "transfer" {
		t.Errorf("Expected transaction type 'transfer', got %s", transaction.TransactionType)
	}

	fromAccount, err := service.GetAccount("test-from-account")
	if err != nil {
		t.Fatalf("Failed to get from account: %v", err)
	}

	toAccount, err := service.GetAccount("test-to-account")
	if err != nil {
		t.Fatalf("Failed to get to account: %v", err)
	}

	expectedFromBalance := 700.0
	expectedToBalance := 300.0

	if fromAccount.Balance != expectedFromBalance {
		t.Errorf("Expected from balance %f, got %f", expectedFromBalance, fromAccount.Balance)
	}

	if toAccount.Balance != expectedToBalance {
		t.Errorf("Expected to balance %f, got %f", expectedToBalance, toAccount.Balance)
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}