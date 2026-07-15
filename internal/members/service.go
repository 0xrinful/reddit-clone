package members

import (
	"context"
	"errors"

	"github.com/0xrinful/reddit-clone/internal/domain"
	"github.com/0xrinful/reddit-clone/internal/shared/errs"
)

type Service interface {
	Join(ctx context.Context, communityID, userID int64) error
	Leave(ctx context.Context, communityID, userID int64) error
	List(ctx context.Context, p ListParams) ([]*MembershipView, error)
	PromoteToModerator(ctx context.Context, actorID int64, p PromoteParams) error
	DemoteModerator(ctx context.Context, actorID int64, p DemoteParams) error
	ListMods(ctx context.Context, p ListParams) ([]*MembershipView, error)
}

type Authorizer interface {
	CanManageModerators(ctx context.Context, communityID, actorID, targetID int64) error
}

type UsersRepo interface {
	GetIDByUsername(ctx context.Context, username string) (int64, error)
}

type service struct {
	authz       Authorizer
	membersRepo Repository
	usersRepo   UsersRepo
}

func NewService(authz Authorizer, membersRepo Repository, usersRepo UsersRepo) Service {
	return &service{authz: authz, membersRepo: membersRepo, usersRepo: usersRepo}
}

func (s *service) Join(ctx context.Context, communityID, userID int64) error {
	return s.membersRepo.Create(ctx, communityID, userID, domain.RoleMember)
}

func (s *service) Leave(ctx context.Context, communityID, userID int64) error {
	membership, err := s.membersRepo.Get(ctx, communityID, userID)
	if errors.Is(err, errs.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}

	if membership.Role.IsOwner() {
		return errs.ErrOwnershipTransferRequired
	}
	return s.membersRepo.Delete(ctx, communityID, userID)
}

func (s *service) List(ctx context.Context, p ListParams) ([]*MembershipView, error) {
	return s.membersRepo.List(ctx, p)
}

func (s *service) PromoteToModerator(ctx context.Context, actorID int64, p PromoteParams) error {
	targetID, err := s.usersRepo.GetIDByUsername(ctx, p.Username)
	if err != nil {
		return err
	}

	if err := s.authz.CanManageModerators(ctx, p.CommunityID, actorID, targetID); err != nil {
		return err
	}

	return s.membersRepo.UpdateAuthority(ctx, p.CommunityID, targetID,
		domain.Authority{
			Role:       domain.RoleModerator,
			Permission: p.Perms,
		},
	)
}

func (s *service) DemoteModerator(ctx context.Context, actorID int64, p DemoteParams) error {
	targetID, err := s.usersRepo.GetIDByUsername(ctx, p.Username)
	if err != nil {
		return err
	}

	if err := s.authz.CanManageModerators(ctx, p.CommunityID, actorID, targetID); err != nil {
		return err
	}

	return s.membersRepo.UpdateAuthority(ctx, p.CommunityID, targetID, domain.Authority{
		Role:       domain.RoleMember,
		Permission: domain.DefaultPerms(domain.RoleMember),
	})
}

func (s *service) ListMods(ctx context.Context, p ListParams) ([]*MembershipView, error) {
	role := domain.RoleModerator
	p.Role = &role
	return s.membersRepo.List(ctx, p)
}
