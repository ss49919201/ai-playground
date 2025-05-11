package auth

import (
	"fmt"
)

// ErrAuthenticationFailure は、認証に失敗したことを表すエラーです。
type ErrAuthenticationFailure struct {
	err error
}

func (e *ErrAuthenticationFailure) Error() string {
	return fmt.Sprintf("failed to authentication: %s", e.err)
}

func NewErrAuthenticationFailure(err error) error {
	return &ErrAuthenticationFailure{err: err}
}
