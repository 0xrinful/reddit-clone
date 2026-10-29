package comments

import (
	"time"

	"github.com/0xrinful/reddit-clone/internal/domain"
	"github.com/0xrinful/reddit-clone/internal/shared/pagination"
)

type Comment struct {
	ID         int64
	PostID     int64
	UserID     int64
	ParentID   *int64
	Content    string
	Score      int64
	ReplyCount int
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

type SortBy string

const (
	SortByNew SortBy = "new"
	SortByTop SortBy = "top"
)

func (s SortBy) Valid() bool {
	switch s {
	case SortByNew, SortByTop:
		return true
	}
	return false
}

type ListByPostParams struct {
	Sort       SortBy
	Pagination pagination.CursorParams[pagination.CommentCursor]
	PostID     int64
}

type ListRepliesParams struct {
	Pagination pagination.CursorParams[pagination.CommentCursor]
	CommentID  int64
}
