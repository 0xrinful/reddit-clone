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

type CreatePostParams struct {
	UserID      int64
	CommunityID int64
	Title       string
	Body        string
}

type UpdatePostParams struct {
	ID          int64
	UserID      int64
	CommunityID int64
	Title       *string
	Body        *string
}
