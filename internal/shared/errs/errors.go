package errs

import "errors"

var (
	ErrEditConflict = errors.New("edit conflict")
	ErrNotFound     = errors.New("not found")
	ErrDuplicate    = errors.New("duplicate")
)
