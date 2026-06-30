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

type PostAuthor struct {
	Username string
}

type PostCommunity struct {
	Name string
}

type PostView struct {
	Post
	Author    PostAuthor
	Community PostCommunity
}

type PostSummary struct {
	ID        int64
	Title     string
	Body      string
	Score     int64
	CreatedAt time.Time
	Author    PostAuthor
	Community PostCommunity
}

type CreateParams struct {
	UserID      int64
	CommunityID int64
	Title       string
	Body        string
}

type UpdateParams struct {
	Title *string
	Body  *string
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

type ListParams struct {
	Sort        SortBy
	Pagination  pagination.Params
	CommunityID int64
}
