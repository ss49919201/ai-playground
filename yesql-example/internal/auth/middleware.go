package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	SessionIDKey contextKey = "session_id"
)

type Middleware struct {
	authService *Service
}

func NewMiddleware(authService *Service) *Middleware {
	return &Middleware{
		authService: authService,
	}
}

func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := m.extractSessionID(r)
		if sessionID == "" {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		session, err := m.authService.ValidateSession(sessionID)
		if err != nil {
			http.Error(w, "Invalid or expired session", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, session.UserID)
		ctx = context.WithValue(ctx, SessionIDKey, session.SessionID)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) extractSessionID(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	cookie, err := r.Cookie("session_id")
	if err == nil {
		return cookie.Value
	}

	return ""
}

func (m *Middleware) RequireAccountAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserIDFromContext(r.Context())
		if userID == "" {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}

		accountID := m.extractAccountID(r)
		if accountID == "" {
			http.Error(w, "Account ID required", http.StatusBadRequest)
			return
		}

		hasAccess, err := m.authService.CheckUserAccountAccess(userID, accountID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to check account access: %v", err), http.StatusInternalServerError)
			return
		}

		if !hasAccess {
			http.Error(w, "Access denied to this account", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) extractAccountID(r *http.Request) string {
	if accountID := r.URL.Query().Get("account_id"); accountID != "" {
		return accountID
	}

	if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
		if err := r.ParseForm(); err == nil {
			if accountID := r.Form.Get("account_id"); accountID != "" {
				return accountID
			}
		}
	}

	return ""
}

func GetUserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}

func GetSessionIDFromContext(ctx context.Context) string {
	sessionID, ok := ctx.Value(SessionIDKey).(string)
	if !ok {
		return ""
	}
	return sessionID
}