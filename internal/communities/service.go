package communities

import (
	"context"

	"github.com/0xrinful/reddit-clone/internal/shared/errs"
)

type Service interface {
	GetByName(ctx context.Context, name string) (*Community, error)
	GetViewByName(ctx context.Context, name string) (*CommunityView, error)
	Create(ctx context.Context, p CreateParams) (*Community, error)
	Delete(ctx context.Context, name string, requesterID int64) error
	Update(ctx context.Context, name string, requesterID int64, p UpdateParams) (*Community, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) GetByName(ctx context.Context, name string) (*Community, error) {
	return s.repo.GetByName(ctx, name)
}

func (s *service) GetViewByName(ctx context.Context, name string) (*CommunityView, error) {
	return s.repo.GetViewByName(ctx, name)
}

func (s *service) Create(ctx context.Context, p CreateParams) (*Community, error) {
	c := &Community{
		Description: p.Description,
		Name:        p.Name,
		OwnerID:     &p.OwnerID,
	}

	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}

	return c, nil
}

func (s *service) Delete(ctx context.Context, name string, requesterID int64) error {
	community, err := s.GetByName(ctx, name)
	if err != nil {
		return err
	}

	if community.OwnerID == nil || *community.OwnerID != requesterID {
		return errs.ErrForbidden
	}

	return s.repo.Delete(ctx, community.ID)
}

func (s *service) Update(
	ctx context.Context,
	name string,
	requesterID int64,
	p UpdateParams,
) (*Community, error) {
	community, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	// TODO: allow mods in the feature
	if community.OwnerID == nil || *community.OwnerID != requesterID {
		return nil, errs.ErrForbidden
	}

	return s.repo.Update(ctx, community.ID, p)
}
