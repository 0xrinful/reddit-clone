package errs

import "errors"

var (
	ErrEditConflict      = errors.New("edit conflict")
	ErrNotFound          = errors.New("not found")
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateUsername = errors.New("duplicate username")
	ErrAlreadyActivated  = errors.New("already activated")
	ErrInvalidToken      = errors.New("invalid or expired token")
)
