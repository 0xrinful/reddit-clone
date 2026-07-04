package errs

import "errors"

var (
	ErrEditConflict              = errors.New("edit conflict")
	ErrNotFound                  = errors.New("not found")
	ErrDuplicateEmail            = errors.New("duplicate email")
	ErrDuplicateUsername         = errors.New("duplicate username")
	ErrDuplicateCommunityName    = errors.New("duplicate community name")
	ErrAlreadyActivated          = errors.New("already activated")
	ErrInvalidToken              = errors.New("invalid or expired token")
	ErrInvalidCredentials        = errors.New("invalid credentials")
	ErrRefreshTokenReuse         = errors.New("refresh token reuse")
	ErrForbidden                 = errors.New("forbidden")
	ErrOwnershipTransferRequired = errors.New(
		"community ownership must be transferred before leaving",
	)
)
