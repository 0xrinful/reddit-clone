package posts

import (
	"context"
)

type Service interface {
	GetDetails(ctx context.Context, id int64) (*PostDetails, error)
	Create(ctx context.Context, params CreatePostParams) (*Post, error)
	Update(ctx context.Context, params UpdatePostParams) error
	Delete(ctx context.Context, id, userID int64) error
	ListSummaries(ctx context.Context, params ListPostParams) ([]*PostSummary, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(
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

func (s *service) Update(ctx context.Context, params UpdatePostParams) error {
	return s.repo.Update(ctx, params)
}

func (s *service) GetDetails(ctx context.Context, id int64) (*PostDetails, error) {
	return s.repo.GetDetails(ctx, id)
}

func (s *service) Delete(ctx context.Context, id, userID int64) error {
	return s.repo.Delete(ctx, id, userID)
}

func (s *service) ListSummaries(
	ctx context.Context,
	params ListPostParams,
) ([]*PostSummary, error) {
	return s.repo.ListSummaries(ctx, params)
}
