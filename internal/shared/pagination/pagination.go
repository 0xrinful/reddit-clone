package pagination

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"time"
)

const (
	DefaultLimit = 25
	MaxLimit     = 100
)

var ErrInvalidCursor = errors.New("invalid cursor")

type Cursor struct {
	ID        int64
	CreatedAt *time.Time
	Score     *int64
}

func (c *Cursor) Encode() string {
	var buf [32]byte // 1 byte mask + up to 3 x 8 byte fields = 25 bytes max, 32 for alignment
	var mask byte
	offset := 1

	mask |= (1 << 0)
	binary.LittleEndian.PutUint64(buf[offset:offset+8], uint64(c.ID))
	offset += 8

	if c.CreatedAt != nil {
		mask |= (1 << 1)
		binary.LittleEndian.PutUint64(buf[offset:offset+8], uint64(c.CreatedAt.Unix()))
		offset += 8
	}

	if c.Score != nil {
		mask |= (1 << 2)
		binary.LittleEndian.PutUint64(buf[offset:offset+8], uint64(*c.Score))
		offset += 8
	}

	buf[0] = mask

	return base64.URLEncoding.EncodeToString(buf[:offset])
}

func Decode(s string) (*Cursor, error) {
	data, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	if len(data) < 9 {
		return nil, ErrInvalidCursor
	}

	var c Cursor
	mask := data[0]
	offset := 1

	c.ID = int64(binary.LittleEndian.Uint64(data[offset : offset+8]))
	offset += 8
	if c.ID < 1 {
		return nil, ErrInvalidCursor
	}

	if mask&(1<<1) != 0 {
		if len(data) < offset+8 {
			return nil, ErrInvalidCursor
		}

		ts := int64(binary.LittleEndian.Uint64(data[offset : offset+8]))
		t := time.Unix(ts, 0).UTC()
		if t.Year() < 2024 || t.Year() > 2100 {
			return nil, ErrInvalidCursor
		}
		c.CreatedAt = &t
		offset += 8
	}

	if mask&(1<<2) != 0 {
		if len(data) < offset+8 {
			return nil, ErrInvalidCursor
		}

		score := int64(binary.LittleEndian.Uint64(data[offset : offset+8]))
		c.Score = &score
	}

	return &c, nil
}

type Params struct {
	Limit  int
	Cursor *Cursor
}
