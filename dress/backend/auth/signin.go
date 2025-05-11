package auth

import (
	"database/sql"

	"github.com/ss49919201/ai-kata/dress/backend/database/mysql"
	"golang.org/x/crypto/bcrypt"
)

type SigninResponse struct {
	Token     string
	ExpiresAt int64
}

// Signin は、ユーザーを認証します。
func Signin(
	email, password string,
) (
	*SigninResponse,
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

	var hashedPassword string
	if err := mysqlClient.QueryRow(
		"SELECT user_password_authentications.password FROM users INNER JOIN user_password_authentications ON users.id = user_password_authentications.user_id WHERE users.email = ?",
		email,
	).Scan(
		&hashedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, NewErrAuthenticationFailure(err)
		}

		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	); err != nil {
		return nil, NewErrAuthenticationFailure(err)
	}

	// Cookie に保存するトークンを生成
	token, expiresAt := generateToken()
	saveToken(token, expiresAt)

	return &SigninResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}
