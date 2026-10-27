package comments

import (
	"time"

	"github.com/0xrinful/reddit-clone/internal/domain"
)

type Comment struct {
	ID         int64
	PostID     int64
	UserID     int64
	ParentID   *int64
	Content    string
	Score      int64
	ReplyCount int
	Depth      int
	Status     domain.CommentStatus
	CreatedAt  time.Time
}

type CommentAuthor struct {
	Username string
}

type CommentView struct {
	Comment
	Author CommentAuthor
}

type CreateParams struct {
	PostID  int64
	UserID  int64
	Content string
}

type ReplyParams struct {
	UserID   int64
	ParentID int64
	Content  string
}
