package users

import (
	"time"

	"github.com/alexedwards/argon2id"
)

type User struct {
	ID          int64
	Username    string
	Email       string
	Password    Password
	AvatarUrl   *string
	CreatedAt   time.Time
	Version     int32
	Activated   bool
	ActivatedAt *time.Time
}

type Password struct {
	hash string
}

func (p *Password) Set(plain string) error {
	hash, err := argon2id.CreateHash(plain, argon2id.DefaultParams)
	if err != nil {
		return err
	}
	p.hash = hash
	return nil
}

func (p *Password) Match(plain string) (bool, error) {
	return argon2id.ComparePasswordAndHash(plain, p.hash)
}
