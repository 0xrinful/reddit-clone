package errs

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrEditConflict = errors.New("edit conflict")
)
