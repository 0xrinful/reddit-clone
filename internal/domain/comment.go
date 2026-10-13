package domain

type CommentStatus string

const (
	CommentStatusVisible CommentStatus = "visible"
	CommentStatusDeleted CommentStatus = "deleted"
	CommentStatusRemoved CommentStatus = "removed"
)

func (s CommentStatus) CanView() bool {
	return s == CommentStatusVisible || s == CommentStatusDeleted
}
