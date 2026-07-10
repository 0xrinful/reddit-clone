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
	ErrPermissionDenied          = errors.New("permission denied")
	ErrForbidden                 = errors.New("forbidden")
	ErrSelfBan                   = errors.New("can't ban self")
	ErrOwnershipTransferRequired = errors.New(
		"community ownership must be transferred before leaving",
	)
)
