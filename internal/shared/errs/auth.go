package errs

import "errors"

var (
	ErrAlreadyActivated   = errors.New("already activated")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid or expired token")
)
