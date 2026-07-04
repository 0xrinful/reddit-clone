package pagination

import (
	"errors"
)

const (
	DefaultLimit = 25
	MaxLimit     = 100
)

var ErrInvalidCursor = errors.New("invalid cursor")

type Params[T any] struct {
	Limit  int
	Cursor *T
}
