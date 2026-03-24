package posts

import (
	"time"
)

type Post struct {
	ID          int64
	Title       string
	Body        string
	UserID      int64
	CommunityID int64
	Views       int64
	CreatedAt   time.Time
	Version     int32
}
