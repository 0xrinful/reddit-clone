package pagination

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type CommunityCursor struct {
	ID        int64      `json:"id"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

func (c *CommunityCursor) Encode() string {
	b, _ := json.Marshal(c)
	return base64.RawURLEncoding.EncodeToString(b)
}

func DecodeCommunityCursor(s string) (*CommunityCursor, error) {
	if s == "" {
		return nil, nil
	}

	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil, ErrInvalidCursor
	}

	var c CommunityCursor
	if err = json.Unmarshal(b, &c); err != nil {
		return nil, ErrInvalidCursor
	}

	return &c, nil
}
