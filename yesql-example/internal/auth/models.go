package auth

import (
	"time"
)

type User struct {
	UserID       string    `json:"user_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Session struct {
	SessionID string    `json:"session_id"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type UserAccount struct {
	UserID      string    `json:"user_id"`
	AccountID   string    `json:"account_id"`
	AccountName string    `json:"account_name"`
	Balance     float64   `json:"balance"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User      *User   `json:"user"`
	SessionID string  `json:"session_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}