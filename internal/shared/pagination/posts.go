package pagination

import (
	"time"
)

type PostCursor struct {
	ID        int64      `json:"id"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	Score     *int64     `json:"score,omitempty"`
}

func (c *PostCursor) Encode() string {
	return EncodeCursor(c)
}

func DecodePostCursor(s string) (*PostCursor, error) {
	return DecodeCursor[PostCursor](s)
}
