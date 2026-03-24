package communities

import "context"

type Service interface {
	GetByName(ctx context.Context, name string) (*Community, error)
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
