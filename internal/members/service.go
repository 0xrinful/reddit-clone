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
}

type service struct {
	membersRepo Repository
}

func NewService(membersRepo Repository) Service {
	return &service{membersRepo: membersRepo}
}

func (s *service) Join(ctx context.Context, communityID, userID int64) error {
	// TODO: check ban state
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
