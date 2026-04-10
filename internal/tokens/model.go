package tokens

import "time"

type Token struct {
	Plaintext string
	Hash      string
	UserID    int64
	Expiry    time.Time
	Scope     string
}
