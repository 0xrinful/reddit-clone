package pagination

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

const (
	DefaultLimit = 25
	MaxLimit     = 100
)

type Cursor struct {
	ID        int64
	CreatedAt *time.Time
	Score     *int64
}

func (c *Cursor) Encode() (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func Decode(s string) (*Cursor, error) {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	var c Cursor
	if err = json.Unmarshal(b, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

type Params struct {
	Limit  int
	Cursor *Cursor
}
