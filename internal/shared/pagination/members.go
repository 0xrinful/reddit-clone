package pagination

import (
	"time"
)

type MemberCursor struct {
	UserID   int64     `json:"user_id"`
	JoinedAt time.Time `json:"joined_at"`
}

func (c *MemberCursor) Encode() string {
	return EncodeCursor(c)
}

func DecodeMemberCursor(s string) (*MemberCursor, error) {
	return DecodeCursor[MemberCursor](s)
}
