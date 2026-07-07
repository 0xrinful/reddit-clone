package pagination

import (
	"encoding/base64"
	"encoding/json"
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

func EncodeCursor[T any](c *T) string {
	b, err := json.Marshal(c)
	if err != nil {
		panic("pagination: invalid cursor")
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

func DecodeCursor[T any](s string) (*T, error) {
	if s == "" {
		return nil, nil
	}

	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil, ErrInvalidCursor
	}

	var c T
	if err = json.Unmarshal(b, &c); err != nil {
		return nil, ErrInvalidCursor
	}

	return &c, nil
}
