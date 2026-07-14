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
