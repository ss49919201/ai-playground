package auth

import (
	"time"

	"github.com/google/uuid"
	"github.com/ss49919201/ai-kata/dress/backend/database/memory/ttlcache"
)

func generateToken() (token string, expiresAt int64) {
	token = uuid.New().String()
	expiresAt = time.Now().Add(time.Hour * 24).Unix()
	return
}

func saveToken(token string, expiresAt int64) {
	ttlcache.Set(token, "", time.Duration(expiresAt))
}

func loadToken(token string) (string, error) {
	return ttlcache.Load(token)
}

func deleteToken(token string) {
	ttlcache.Delete(token)
}
