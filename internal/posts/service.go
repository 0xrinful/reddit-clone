package posts

import "context"

type Service interface {
	GetPost(ctx context.Context, id int64) (*Post, error)
	CreatePost(ctx context.Context, authorID int64, req CreatePostRequest) (*Post, error)
}

type service struct {
	repo Repository
	// validator *validator.Validate
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) CreatePost(
	ctx context.Context,
	authorID int64,
	req CreatePostRequest,
) (*Post, error) {
	return nil, nil
}

func (s *service) GetPost(ctx context.Context, id int64) (*Post, error) {
	return nil, nil
}
