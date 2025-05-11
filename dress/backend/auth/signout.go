package auth

func Signout(token string) error {
	deleteToken(token)
	return nil
}
