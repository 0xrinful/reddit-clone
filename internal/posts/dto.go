package posts

import (
	"time"

	"github.com/0xrinful/reddit-clone/internal/shared/validator"
)

// validation helpers
const (
	postTitleMin = 3
	postTitleMax = 120
	postBodyMin  = 10
	postBodyMax  = 40000
)

func validateTitle(v *validator.Validator, title string) {
	v.Check(validator.NotBlank(title), "title", "must not be blank")
	v.Check(validator.MinLength(title, postTitleMin), "title", "must be at least 3 characters")
	v.Check(validator.MaxLength(title, postTitleMax), "title", "must not exceed 120 characters")
}

func validateBody(v *validator.Validator, body string) {
	v.Check(validator.NotBlank(body), "body", "must not be blank")
	v.Check(validator.MinLength(body, postBodyMin), "body", "must be at least 10 characters")
	v.Check(validator.MaxLength(body, postBodyMax), "body", "must not exceed 40000 characters")
}

// request structs
type CreatePostRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (r *CreatePostRequest) Validate(v *validator.Validator) {
	validateTitle(v, r.Title)
	validateBody(v, r.Body)
}

type UpdatePostRequest struct {
	Title *string `json:"title"`
	Body  *string `json:"body"`
}

func (r *UpdatePostRequest) Validate(v *validator.Validator) {
	if r.Title == nil && r.Body == nil {
		v.AddError("request", "must provide at least one field")
		return
	}
	if r.Title != nil {
		validateTitle(v, *r.Title)
	}
	if r.Body != nil {
		validateBody(v, *r.Body)
	}
}

// DTOs
type PostPublicDTO struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Body          string    `json:"body"`
	Score         int64     `json:"score"`
	UserID        int64     `json:"user_id"`
	CommunityName string    `json:"community"`
	CreatedAt     time.Time `json:"created_at"`
}

type PostOwnerDTO struct {
	PostPublicDTO
	Views   int64 `json:"views"`
	Version int32 `json:"version"`
}

// mapping helpers
func toPostPublicDTO(p *Post, communityName string) PostPublicDTO {
	return PostPublicDTO{
		ID:            p.ID,
		Title:         p.Title,
		Body:          p.Body,
		UserID:        p.UserID,
		Score:         p.Score,
		CommunityName: communityName,
		CreatedAt:     p.CreatedAt,
	}
}

func toPostOwnerDTO(p *Post, communityName string) PostOwnerDTO {
	return PostOwnerDTO{
		PostPublicDTO: toPostPublicDTO(p, communityName),
		Views:         p.Views,
		Version:       p.Version,
	}
}

// response envelope
type PostResponse struct {
	Post PostPublicDTO `json:"post"`
}

type PostOwnerResponse struct {
	Post PostOwnerDTO `json:"post"`
}

type ListPostsResponse struct {
	Posts      []PostPublicDTO `json:"posts"`
	NextCursor string          `json:"next_cursor,omitempty"`
}

// response constructor
func toPostResponse(p *Post, communityName string) PostResponse {
	return PostResponse{
		Post: toPostPublicDTO(p, communityName),
	}
}

func toPostOwnerResponse(p *Post, communityName string) PostOwnerResponse {
	return PostOwnerResponse{
		Post: toPostOwnerDTO(p, communityName),
	}
}

func toListPostsResponse(p []*Post, nextCursor string, communityName string) ListPostsResponse {
	posts := make([]PostPublicDTO, len(p))
	for i := range p {
		posts[i] = toPostPublicDTO(p[i], communityName)
	}
	return ListPostsResponse{Posts: posts, NextCursor: nextCursor}
}
