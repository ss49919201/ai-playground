package auth

import (
	"database/sql"
	"os"
	"testing"
	"yesql-account-system/internal/db"
	"yesql-account-system/internal/yesql"

	_ "github.com/mattn/go-sqlite3"
)

func setupAuthTestDB(t *testing.T) (*sql.DB, *Service) {
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

func TestRegister(t *testing.T) {
	dbConn, service := setupAuthTestDB(t)
	defer dbConn.Close()

	req := RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	user, err := service.Register(req)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	if user.Username != req.Username {
		t.Errorf("Expected username %s, got %s", req.Username, user.Username)
	}

	if user.Email != req.Email {
		t.Errorf("Expected email %s, got %s", req.Email, user.Email)
	}

	if user.UserID == "" {
		t.Error("Expected user ID to be set")
	}
}

func TestLogin(t *testing.T) {
	dbConn, service := setupAuthTestDB(t)
	defer dbConn.Close()

	registerReq := RegisterRequest{
		Username: "testuser2",
		Email:    "test2@example.com",
		Password: "password123",
	}

	_, err := service.Register(registerReq)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	loginReq := LoginRequest{
		Username: "testuser2",
		Password: "password123",
	}

	loginResponse, err := service.Login(loginReq)
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	if loginResponse.User.Username != registerReq.Username {
		t.Errorf("Expected username %s, got %s", registerReq.Username, loginResponse.User.Username)
	}

	if loginResponse.SessionID == "" {
		t.Error("Expected session ID to be set")
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	dbConn, service := setupAuthTestDB(t)
	defer dbConn.Close()

	registerReq := RegisterRequest{
		Username: "testuser3",
		Email:    "test3@example.com",
		Password: "password123",
	}

	_, err := service.Register(registerReq)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	loginReq := LoginRequest{
		Username: "testuser3",
		Password: "wrongpassword",
	}

	_, err = service.Login(loginReq)
	if err == nil {
		t.Error("Expected error for invalid password, but got none")
	}
}

func TestValidateSession(t *testing.T) {
	dbConn, service := setupAuthTestDB(t)
	defer dbConn.Close()

	registerReq := RegisterRequest{
		Username: "testuser4",
		Email:    "test4@example.com",
		Password: "password123",
	}

	user, err := service.Register(registerReq)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	session, err := service.CreateSession(user.UserID)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	validatedSession, err := service.ValidateSession(session.SessionID)
	if err != nil {
		t.Fatalf("Failed to validate session: %v", err)
	}

	if validatedSession.UserID != user.UserID {
		t.Errorf("Expected user ID %s, got %s", user.UserID, validatedSession.UserID)
	}
}

func TestValidateInvalidSession(t *testing.T) {
	dbConn, service := setupAuthTestDB(t)
	defer dbConn.Close()

	_, err := service.ValidateSession("invalid-session-id")
	if err == nil {
		t.Error("Expected error for invalid session, but got none")
	}
}

func TestLogout(t *testing.T) {
	dbConn, service := setupAuthTestDB(t)
	defer dbConn.Close()

	registerReq := RegisterRequest{
		Username: "testuser5",
		Email:    "test5@example.com",
		Password: "password123",
	}

	user, err := service.Register(registerReq)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	session, err := service.CreateSession(user.UserID)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	err = service.Logout(session.SessionID)
	if err != nil {
		t.Fatalf("Failed to logout: %v", err)
	}

	_, err = service.ValidateSession(session.SessionID)
	if err == nil {
		t.Error("Expected session to be invalid after logout, but it's still valid")
	}
}

func TestCreateUserAccountAssociation(t *testing.T) {
	dbConn, service := setupAuthTestDB(t)
	defer dbConn.Close()

	registerReq := RegisterRequest{
		Username: "testuser6",
		Email:    "test6@example.com",
		Password: "password123",
	}

	user, err := service.Register(registerReq)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	accountID := "test-account-1"
	
	// First create an account in the accounts table
	_, err = dbConn.Exec("INSERT INTO accounts (account_id, account_name, balance) VALUES (?, ?, ?)", 
		accountID, "Test Account", 1000.0)
	if err != nil {
		t.Fatalf("Failed to create test account: %v", err)
	}

	err = service.CreateUserAccountAssociation(user.UserID, accountID)
	if err != nil {
		t.Fatalf("Failed to create user account association: %v", err)
	}

	hasAccess, err := service.CheckUserAccountAccess(user.UserID, accountID)
	if err != nil {
		t.Fatalf("Failed to check user account access: %v", err)
	}

	if !hasAccess {
		t.Error("Expected user to have access to account")
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}