package posts

import "github.com/0xrinful/reddit-clone/internal/shared/validator"

// request structs
type CreatePostRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (r *CreatePostRequest) Validate(v *validator.Validator) {
	v.Check(validator.NotBlank(r.Title), "title", "must not be blank")
	v.Check(validator.MaxLength(r.Title, 100), "title", "must not exceed 100 characters")
	v.Check(validator.NotBlank(r.Body), "body", "must not be blank")
	v.Check(validator.MaxLength(r.Body, 40000), "body", "must not exceed 40000 characters")
}

type UpdatePostRequest struct {
	Title *string `json:"title"`
	Body  *string `json:"body"`
}

// DTOs
type PostPublicDTO struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type PostOwnerDTO struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
	// extra fields for owner
}

// mapping helpers
func toPostPublicDTO(p *Post) PostPublicDTO {
	return PostPublicDTO{
		ID:    p.ID,
		Title: p.Title,
		Body:  p.Body,
	}
}

func toPostOwnerDTO(p *Post) PostOwnerDTO {
	return PostOwnerDTO{
		ID:    p.ID,
		Title: p.Title,
		Body:  p.Body,
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
