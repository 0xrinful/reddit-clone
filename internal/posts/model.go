package posts

import (
	"time"

	"github.com/0xrinful/reddit-clone/internal/shared/pagination"
)

type Post struct {
	ID          int64
	Title       string
	Body        string
	UserID      int64
	CommunityID int64
	Views       int64
	Score       int64
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

type SortBy string

const (
	SortByNew SortBy = "new"
	SortByTop SortBy = "top"
	SortByHot SortBy = "hot" // TODO: implement hot sort later
)

func (s SortBy) IsValid() bool {
	switch s {
	case SortByNew, SortByTop, SortByHot:
		return true
	default:
		return false
	}
}

func (s SortBy) ToSql() string {
	switch s {
	case SortByNew:
		return "created_at DESC"
	case SortByTop:
		return "score DESC"
	default:
		return "created_at DESC"
	}
}

type ListPostParams struct {
	Sort        SortBy
	Cursor      pagination.Cursor
	CommunityID int64
}
