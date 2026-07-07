package pagination

import (
	"time"
)

type CommunityCursor struct {
	ID        int64      `json:"id"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

func (c *CommunityCursor) Encode() string {
	return EncodeCursor(c)
}

func DecodeCommunityCursor(s string) (*CommunityCursor, error) {
	return DecodeCursor[CommunityCursor](s)
}
