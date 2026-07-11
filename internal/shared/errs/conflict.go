package errs

import "errors"

var (
	ErrDuplicateEmail         = errors.New("duplicate email")
	ErrDuplicateUsername      = errors.New("duplicate username")
	ErrDuplicateCommunityName = errors.New("duplicate community name")
)
