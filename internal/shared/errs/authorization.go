package errs

import "errors"

var (
	ErrPermissionDenied = errors.New("permission denied")
	ErrSelfBan          = errors.New("can't ban self")
)
