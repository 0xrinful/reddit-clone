package posts

import (
	"context"

	"github.com/0xrinful/reddit-clone/internal/shared/errs"
)

type Service interface {
	Get(ctx context.Context, id int64) (*PostView, error)
	Create(ctx context.Context, params CreateParams) (*Post, error)
	Update(ctx context.Context, id, requesterID int64, p UpdateParams) (*Post, error)
	Delete(ctx context.Context, id, userID int64) error
	List(ctx context.Context, params ListParams) ([]*PostSummary, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(
	ctx context.Context,
	params CreateParams,
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

func (s *service) Update(
	ctx context.Context,
	id, requesterID int64,
	p UpdateParams,
) (*Post, error) {
	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if post.UserID != requesterID {
		return nil, errs.ErrForbidden
	}

	return s.repo.Update(ctx, id, p)
}

func (s *service) Get(ctx context.Context, id int64) (*PostView, error) {
	return s.repo.GetView(ctx, id)
}

func (s *service) Delete(ctx context.Context, id, requesterID int64) error {
	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if post.UserID != requesterID {
		return errs.ErrForbidden
	}

	return s.repo.Delete(ctx, id)
}

func (s *service) List(
	ctx context.Context,
	params ListParams,
) ([]*PostSummary, error) {
	return s.repo.ListSummaries(ctx, params)
}
