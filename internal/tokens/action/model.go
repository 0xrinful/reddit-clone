package action

import (
	"time"

	"github.com/0xrinful/reddit-clone/internal/tokens"
)

const (
	ScopeActivation    = "activation"
	ScopePasswordReset = "password-reset"
)

type Token struct {
	Plaintext string
	Hash      []byte
	UserID    int64
	Expiry    time.Time
	Scope     string
}

func Generate(userID int64, ttl time.Duration, scope string) (*Token, error) {
	plain, hash, err := tokens.GenerateToken()
	if err != nil {
		return nil, err
	}

	return &Token{
		Plaintext: plain,
		Hash:      hash,
		UserID:    userID,
		Expiry:    time.Now().Add(ttl).UTC(),
		Scope:     scope,
	}, nil
}
