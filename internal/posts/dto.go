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
	validator.ValidatePostTitle(v, r.Title)
	validator.ValidatePostBody(v, r.Body)
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
		validator.ValidatePostTitle(v, *r.Title)
	}
	if r.Body != nil {
		validator.ValidatePostBody(v, *r.Body)
	}
}

// DTOs
type AuthorDTO struct {
	ID       int64  `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}

type CommunityDTO struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type PostPublicDTO struct {
	ID        int64        `json:"id"`
	Title     string       `json:"title"`
	Body      string       `json:"body"`
	Score     int64        `json:"score"`
	Author    AuthorDTO    `json:"author"`
	Community CommunityDTO `json:"community"`
	CreatedAt time.Time    `json:"created_at"`
}

type PostOwnerDTO struct {
	PostPublicDTO
	Views int64 `json:"views"`
}

type PostSummaryDTO struct {
	ID        int64        `json:"id"`
	Title     string       `json:"title"`
	Body      string       `json:"body"`
	Score     int64        `json:"score"`
	CreatedAt time.Time    `json:"created_at"`
	Author    AuthorDTO    `json:"author"`
	Community CommunityDTO `json:"community"`
}

// mapping helpers
func postToView(p *Post, communityName string) *PostView {
	return &PostView{
		Post: *p,
		Community: PostCommunity{
			Name: communityName,
		},
	}
}

func toPostPublicDTO(p *PostView) PostPublicDTO {
	return PostPublicDTO{
		ID:    p.ID,
		Title: p.Title,
		Body:  p.Body,
		Score: p.Score,
		Author: AuthorDTO{
			ID:       p.UserID,
			Username: p.Author.Username,
		},
		Community: CommunityDTO{
			ID:   p.CommunityID,
			Name: p.Community.Name,
		},
		CreatedAt: p.CreatedAt,
	}
}

func toPostOwnerDTO(p *PostView) PostOwnerDTO {
	return PostOwnerDTO{
		PostPublicDTO: toPostPublicDTO(p),
		Views:         p.Views,
	}
}

func toPostSummaryDTO(p *PostSummary) PostSummaryDTO {
	return PostSummaryDTO{
		ID:    p.ID,
		Title: p.Title,
		Body:  p.Body,
		Score: p.Score,
		Author: AuthorDTO{
			Username: p.Author.Username,
		},
		Community: CommunityDTO{
			Name: p.Community.Name,
		},
		CreatedAt: p.CreatedAt,
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
	Posts      []PostSummaryDTO `json:"posts"`
	NextCursor string           `json:"next_cursor,omitempty"`
}

// response constructor
func toPostResponse(p *PostView) PostResponse {
	return PostResponse{
		Post: toPostPublicDTO(p),
	}
}

func toPostOwnerResponse(p *PostView) PostOwnerResponse {
	return PostOwnerResponse{
		Post: toPostOwnerDTO(p),
	}
}

func toListPostsResponse(p []*PostSummary, nextCursor string) ListPostsResponse {
	posts := make([]PostSummaryDTO, len(p))
	for i := range p {
		posts[i] = toPostSummaryDTO(p[i])
	}
	return ListPostsResponse{Posts: posts, NextCursor: nextCursor}
}
