package errs

import "errors"

var (
	ErrPermissionDenied = errors.New("permission denied")
	ErrBanned           = errors.New("banned from community")
	ErrNotMember        = errors.New("not a community member")
	ErrBlocked          = errors.New("blocked by user")
	ErrSelfBan          = errors.New("can't ban self")
)
