package pagination

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type PostCursor struct {
	ID        int64      `json:"id"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	Score     *int64     `json:"score,omitempty"`
}

func (c *PostCursor) Encode() string {
	b, _ := json.Marshal(c)
	return base64.RawURLEncoding.EncodeToString(b)
}

func DecodePostCursor(s string) (*PostCursor, error) {
	if s == "" {
		return nil, nil
	}

	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil, ErrInvalidCursor
	}

	var c PostCursor
	if err = json.Unmarshal(b, &c); err != nil {
		return nil, ErrInvalidCursor
	}

	return &c, nil
}
