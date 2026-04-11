package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

func GenerateToken() (plaintext string, hash []byte, err error) {
	b := make([]byte, 16)
	if _, err = rand.Read(b); err != nil {
		return
	}
	plaintext = base64.RawURLEncoding.EncodeToString(b)
	hash = Hash(plaintext)
	return
}

func Hash(plaintext string) []byte {
	hash := sha256.Sum256([]byte(plaintext))
	return hash[:]
}
