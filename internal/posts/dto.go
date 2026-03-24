package posts

import (
	"time"

	"github.com/0xrinful/reddit-clone/internal/shared/validator"
)

// request structs
type CreatePostRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (r *CreatePostRequest) Validate(v *validator.Validator) {
	v.Check(validator.NotBlank(r.Title), "title", "must not be blank")
	v.Check(validator.MaxLength(r.Title, 120), "title", "must not exceed 120 characters")
	v.Check(validator.NotBlank(r.Body), "body", "must not be blank")
	v.Check(validator.MaxLength(r.Body, 40000), "body", "must not exceed 40000 characters")
}

type UpdatePostRequest struct {
	Title *string `json:"title"`
	Body  *string `json:"body"`
}

// DTOs
type PostPublicDTO struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	UserID      int64     `json:"user_id"`
	CommunityID int64     `json:"community_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type PostOwnerDTO struct {
	PostPublicDTO
	Views   int64 `json:"views"`
	Version int32 `json:"version"`
}

// mapping helpers
func toPostPublicDTO(p *Post) PostPublicDTO {
	return PostPublicDTO{
		ID:          p.ID,
		Title:       p.Title,
		Body:        p.Body,
		UserID:      p.UserID,
		CommunityID: p.CommunityID,
		CreatedAt:   p.CreatedAt,
	}
}

func toPostOwnerDTO(p *Post) PostOwnerDTO {
	return PostOwnerDTO{
		PostPublicDTO: PostPublicDTO{
			ID:          p.ID,
			Title:       p.Title,
			Body:        p.Body,
			UserID:      p.UserID,
			CommunityID: p.CommunityID,
			CreatedAt:   p.CreatedAt,
		},
		Views:   p.Views,
		Version: p.Version,
	}
}

// response envelope
type PostResponse struct {
	Post PostPublicDTO `json:"post"`
}

// response constructor
func toPostResponse(p *Post) PostResponse {
	return PostResponse{
		Post: toPostPublicDTO(p),
	}
}
