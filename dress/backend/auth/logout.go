package auth

func Logout(token string) error {
	deleteToken(token)
	return nil
}
