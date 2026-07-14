package pagination

import "time"

type BanCursor struct {
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *BanCursor) Encode() string {
	return EncodeCursor(c)
}

func DecodeBanCursor(s string) (*BanCursor, error) {
	return DecodeCursor[BanCursor](s)
}
