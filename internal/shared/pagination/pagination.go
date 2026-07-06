package pagination

import (
	"errors"
)

const (
	DefaultLimit = 25
	MaxLimit     = 100
)

var ErrInvalidCursor = errors.New("invalid cursor")

type CursorParams[T any] struct {
	Limit  int
	Cursor *T
}

type OffsetParams struct {
	Limit int
	Page  int
}

func (p OffsetParams) Offset() int {
	return (p.Page - 1) * p.Limit
}
