package comments

import (
	"time"

	"github.com/0xrinful/reddit-clone/internal/shared/validator"
)

// request structs
type CreateCommentRequest struct {
	Content string `json:"content"`
}

func (r *CreateCommentRequest) Vaildate(v *validator.Validator) {
	validator.ValidateCommentContent(v, r.Content)
}

type ReplyCommentRequest struct {
	Content string `json:"content"`
}

func (r *ReplyCommentRequest) Vaildate(v *validator.Validator) {
	validator.ValidateCommentContent(v, r.Content)
}

// DTOs
type AuthorDTO struct {
	ID       int64  `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}

type PostDTO struct {
	ID int64 `json:"id,omitempty"`
}

type CommentPublicDTO struct {
	ID         int64     `json:"id"`
	Post       PostDTO   `json:"post"`
	Author     AuthorDTO `json:"author"`
	ParentID   *int64    `json:"parent_id"`
	Content    string    `json:"content"`
	Score      int64     `json:"score"`
	ReplyCount int       `json:"reply_count"`
	CreatedAt  time.Time `json:"created_at"`
}

// mapping helpers
func commentToView(c *Comment) *CommentView {
	return &CommentView{
		Comment: *c,
	}
}

func toCommentPublicDTO(c *CommentView) CommentPublicDTO {
	return CommentPublicDTO{
		ID: c.ID,
		Post: PostDTO{
			ID: c.PostID,
		},
		Author: AuthorDTO{
			ID: c.UserID,
		},
		ParentID:   c.ParentID,
		Content:    c.Content,
		Score:      c.Score,
		ReplyCount: c.ReplyCount,
		CreatedAt:  c.CreatedAt,
	}
}

// response envelope
type CommentResponse struct {
	Comment any `json:"comment"`
}

// response constructor
func toCommentResponse(c *CommentView) CommentResponse {
	return CommentResponse{Comment: toCommentPublicDTO(c)}
}
