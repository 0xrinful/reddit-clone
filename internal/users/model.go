package users

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
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
	plain *string
	hash  []byte
}

func (p *Password) Set(plain string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), 12)
	if err != nil {
		return err
	}

	p.plain = &plain
	p.hash = hash
	return nil
}

func (p *Password) Match(plain string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plain))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

type CreateUserParams struct {
	Username      string
	Email         string
	PlainPassword string
}
