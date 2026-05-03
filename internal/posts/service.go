package posts

import (
	"context"
)

type Service interface {
	GetPost(ctx context.Context, id int64) (*PostDetails, error)
	CreatePost(ctx context.Context, params CreatePostParams) (*Post, error)
	UpdatePost(ctx context.Context, params UpdatePostParams) error
	DeletePost(ctx context.Context, id, userID int64) error
	List(ctx context.Context, params ListPostParams) ([]*PostSummary, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) CreatePost(
	ctx context.Context,
	params CreatePostParams,
) (*Post, error) {
	p := &Post{
		Title:       params.Title,
		Body:        params.Body,
		UserID:      params.UserID,
		CommunityID: params.CommunityID,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *service) UpdatePost(ctx context.Context, params UpdatePostParams) error {
	return s.repo.Update(ctx, params)
}

func (s *service) GetPost(ctx context.Context, id int64) (*PostDetails, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) DeletePost(ctx context.Context, id, userID int64) error {
	return s.repo.Delete(ctx, id, userID)
}

func (s *service) List(ctx context.Context, params ListPostParams) ([]*PostSummary, error) {
	return s.repo.List(ctx, params)
}
