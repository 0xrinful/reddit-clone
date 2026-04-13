package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"
)

const (
	ScopeActivation    = "activation"
	ScopeAuth          = "auth"
	ScopePasswordReset = "password-reset"
)

type Token struct {
	Plaintext string
	ID        int64
	Hash      []byte
	UserID    int64
	Expiry    time.Time
	Scope     string
}

func Hash(plaintext string) []byte {
	hash := sha256.Sum256([]byte(plaintext))
	return hash[:]
}

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
