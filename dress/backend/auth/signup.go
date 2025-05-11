package auth

import (
	"github.com/ss49919201/ai-kata/dress/backend/database/mysql"
	"golang.org/x/crypto/bcrypt"
)

type SignupResponse struct {
	Token     string
	ExpiresAt int64
}

// Signup は、ユーザーを新規登録します。
func Signup(
	email, password string,
) (
	*SignupResponse,
	error,
) {
	mysqlClient, err := mysql.NewClient(
		"root",
		"password",
		"localhost:3306",
		"dress",
	)
	if err != nil {
		return nil, err
	}

	defer mysqlClient.Close()

	result, err := mysqlClient.Exec(
		"INSERT INTO users (email, password) VALUES (?, ?)",
		email,
		password,
	)
	if err != nil {
		return nil, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	if _, err := mysqlClient.Exec(
		"INSERT INTO user_password_authentications (user_id, password) VALUES (?, ?)",
		lastInsertId,
		hashedPassword,
	); err != nil {
		return nil, err
	}

	token, expiresAt := generateToken()
	saveToken(token, expiresAt)

	return &SignupResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}
