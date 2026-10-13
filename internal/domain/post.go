package domain

type PostStatus string

const (
	PostStatusVisible PostStatus = "visible"
	PostStatusDeleted PostStatus = "deleted"
	PostStatusRemoved PostStatus = "removed"
)

func (s PostStatus) CanView() bool {
	return s == PostStatusVisible || s == PostStatusDeleted
}

type PostAuthorizationData struct {
	AuthorID    int64
	CommunityID int64
	Status      PostStatus
}
