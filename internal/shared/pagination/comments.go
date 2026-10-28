package pagination

import "time"

type CommentCursor struct {
	ID        int64      `json:"id"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	Score     *int64     `json:"score,omitempty"`
}

func (c *CommentCursor) Encode() string {
	return EncodeCursor(c)
}

func DecodeCommentCursor(s string) (*CommentCursor, error) {
	return DecodeCursor[CommentCursor](s)
}
