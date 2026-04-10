package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"
)

func Generate(userID int64, ttl time.Duration, scope string) (*Token, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	plain := base64.RawURLEncoding.EncodeToString(b)
	hash := Hash(plain)

	return &Token{
		Plaintext: plain,
		Hash:      hash,
		UserID:    userID,
		Expiry:    time.Now().Add(ttl).UTC(),
		Scope:     scope,
	}, nil
}

func Hash(plaintext string) string {
	hash := sha256.Sum256([]byte(plaintext))
	return hex.EncodeToString(hash[:])
}
