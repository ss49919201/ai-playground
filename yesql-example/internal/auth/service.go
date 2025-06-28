package auth

import (
	"database/sql"
	"fmt"
	"time"
	"yesql-account-system/internal/yesql"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func (s *Service) Register(req RegisterRequest) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	userID := uuid.New().String()
	
	query, err := s.loader.GetQuery("create_user")
	if err != nil {
		return nil, fmt.Errorf("failed to get create_user query: %w", err)
	}

	_, err = s.db.Exec(query, userID, req.Username, req.Email, string(hashedPassword))
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return s.GetUserByID(userID)
}

func (s *Service) Login(req LoginRequest) (*LoginResponse, error) {
	user, err := s.GetUserByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	session, err := s.CreateSession(user.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &LoginResponse{
		User:      user,
		SessionID: session.SessionID,
		ExpiresAt: session.ExpiresAt,
	}, nil
}

func (s *Service) Logout(sessionID string) error {
	query, err := s.loader.GetQuery("delete_session")
	if err != nil {
		return fmt.Errorf("failed to get delete_session query: %w", err)
	}

	_, err = s.db.Exec(query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (s *Service) GetUserByID(userID string) (*User, error) {
	query, err := s.loader.GetQuery("get_user_by_id")
	if err != nil {
		return nil, fmt.Errorf("failed to get get_user_by_id query: %w", err)
	}

	var user User
	err = s.db.QueryRow(query, userID).Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (s *Service) GetUserByUsername(username string) (*User, error) {
	query, err := s.loader.GetQuery("get_user_by_username")
	if err != nil {
		return nil, fmt.Errorf("failed to get get_user_by_username query: %w", err)
	}

	var user User
	err = s.db.QueryRow(query, username).Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (s *Service) CreateSession(userID string) (*Session, error) {
	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour) // 24時間で期限切れ

	query, err := s.loader.GetQuery("create_session")
	if err != nil {
		return nil, fmt.Errorf("failed to get create_session query: %w", err)
	}

	_, err = s.db.Exec(query, sessionID, userID, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &Session{
		SessionID: sessionID,
		UserID:    userID,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}, nil
}

func (s *Service) ValidateSession(sessionID string) (*Session, error) {
	query, err := s.loader.GetQuery("get_valid_session")
	if err != nil {
		return nil, fmt.Errorf("failed to get get_valid_session query: %w", err)
	}

	var session Session
	err = s.db.QueryRow(query, sessionID).Scan(
		&session.SessionID,
		&session.UserID,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired session: %w", err)
	}

	return &session, nil
}

func (s *Service) CreateUserAccountAssociation(userID, accountID string) error {
	query, err := s.loader.GetQuery("create_user_account_association")
	if err != nil {
		return fmt.Errorf("failed to get create_user_account_association query: %w", err)
	}

	_, err = s.db.Exec(query, userID, accountID)
	if err != nil {
		return fmt.Errorf("failed to create user account association: %w", err)
	}

	return nil
}

func (s *Service) CheckUserAccountAccess(userID, accountID string) (bool, error) {
	query, err := s.loader.GetQuery("check_user_account_access")
	if err != nil {
		return false, fmt.Errorf("failed to get check_user_account_access query: %w", err)
	}

	var count int
	err = s.db.QueryRow(query, userID, accountID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user account access: %w", err)
	}

	return count > 0, nil
}

func (s *Service) GetUserAccounts(userID string) ([]UserAccount, error) {
	query, err := s.loader.GetQuery("get_user_accounts")
	if err != nil {
		return nil, fmt.Errorf("failed to get get_user_accounts query: %w", err)
	}

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user accounts: %w", err)
	}
	defer rows.Close()

	var accounts []UserAccount
	for rows.Next() {
		var account UserAccount
		err := rows.Scan(
			&account.UserID,
			&account.AccountID,
			&account.AccountName,
			&account.Balance,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user account: %w", err)
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *Service) CleanupExpiredSessions() error {
	query, err := s.loader.GetQuery("delete_expired_sessions")
	if err != nil {
		return fmt.Errorf("failed to get delete_expired_sessions query: %w", err)
	}

	_, err = s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}

	return nil
}