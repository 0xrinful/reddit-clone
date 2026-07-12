package communities

import (
	"context"

	"github.com/0xrinful/reddit-clone/internal/database"
	"github.com/0xrinful/reddit-clone/internal/domain"
	"github.com/0xrinful/reddit-clone/internal/shared/errs"
)

type Service interface {
	GetByName(ctx context.Context, name string) (*Community, error)
	GetViewByName(ctx context.Context, name string) (*CommunityView, error)
	Create(ctx context.Context, p CreateParams) (*Community, error)
	Delete(ctx context.Context, name string, actorID int64) error
	Update(ctx context.Context, name string, actorID int64, p UpdateParams) (*Community, error)
	Search(ctx context.Context, p SearchParams) ([]*CommunitySummary, error)
	List(ctx context.Context, p ListParams) ([]*CommunitySummary, error)
}

type MembersRepo interface {
	Create(ctx context.Context, communityID, userID int64, role domain.Role) error
}

type Authorizer interface {
	CanUpdateCommunity(ctx context.Context, communityID, actorID int64) error
	CanDeleteCommunity(ctx context.Context, communityID, actorID int64) error
}

type service struct {
	txBeginner      database.TxBeginner
	authz           Authorizer
	communitiesRepo Repository
	membersRepo     MembersRepo
}

func NewService(
	txBeginner database.TxBeginner,
	authz Authorizer,
	communitiesRepo Repository,
	membersRepo MembersRepo,
) Service {
	return &service{
		txBeginner:      txBeginner,
		authz:           authz,
		communitiesRepo: communitiesRepo,
		membersRepo:     membersRepo,
	}
}

func (s *service) GetByName(ctx context.Context, name string) (*Community, error) {
	return s.communitiesRepo.GetByName(ctx, name)
}

func (s *service) GetViewByName(ctx context.Context, name string) (*CommunityView, error) {
	return s.communitiesRepo.GetViewByName(ctx, name)
}

func (s *service) Create(ctx context.Context, p CreateParams) (*Community, error) {
	c := &Community{
		Description: p.Description,
		Name:        p.Name,
		OwnerID:     &p.OwnerID,
	}

	tx, err := s.txBeginner.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	ctxTx := database.WithTx(ctx, tx)

	if err := s.communitiesRepo.Create(ctxTx, c); err != nil {
		return nil, err
	}

	if err := s.membersRepo.Create(ctxTx, c.ID, p.OwnerID, domain.RoleOwner); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return c, nil
}

func (s *service) Delete(ctx context.Context, name string, actorID int64) error {
	community, err := s.GetByName(ctx, name)
	if err != nil {
		return err
	}

	if err := s.authz.CanDeleteCommunity(ctx, community.ID, actorID); err != nil {
		return err
	}

	return s.communitiesRepo.Delete(ctx, community.ID)
}

func (s *service) Update(
	ctx context.Context,
	name string,
	actorID int64,
	p UpdateParams,
) (*Community, error) {
	community, err := s.communitiesRepo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	if err := s.authz.CanUpdateCommunity(ctx, community.ID, actorID); err != nil {
		return nil, err
	}

	return s.communitiesRepo.Update(ctx, community.ID, p)
}

func (s *service) Search(ctx context.Context, p SearchParams) ([]*CommunitySummary, error) {
	return s.communitiesRepo.Search(ctx, p)
}

func (s *service) List(ctx context.Context, p ListParams) ([]*CommunitySummary, error) {
	return s.communitiesRepo.List(ctx, p)
}
