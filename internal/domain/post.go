package domain

type PostStatus string

const (
	PostStatusActive       PostStatus = "active"
	PostStatusRemovedByMod PostStatus = "removed_by_mod"
	PostStatusBanned       PostStatus = "banned"
)

func (s PostStatus) IsVisible() bool {
	return s == PostStatusActive
}

type PostAuthzInfo struct {
	AuthorID    int64
	CommunityID int64
	Status      PostStatus
}
