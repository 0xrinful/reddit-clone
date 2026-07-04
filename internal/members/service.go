package members

import (
	"context"
	"errors"

	"github.com/0xrinful/reddit-clone/internal/shared/errs"
)

type Service interface {
	Join(ctx context.Context, communityID, userID int64) error
	Leave(ctx context.Context, communityID, userID int64) error
}

type service struct {
	membersRepo Repository
}

func NewService(membersRepo Repository) Service {
	return &service{membersRepo: membersRepo}
}

func (s *service) Join(ctx context.Context, communityID, userID int64) error {
	// TODO: check ban state
	m := &Membership{
		CommunityID: communityID,
		UserID:      userID,
		Role:        RoleMember,
	}
	return s.membersRepo.Create(ctx, m)
}

func (s *service) Leave(ctx context.Context, communityID, userID int64) error {
	membership, err := s.membersRepo.Get(ctx, communityID, userID)
	if errors.Is(err, errs.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}

	if membership.Role == RoleOwner {
		return errs.ErrOwnershipTransferRequired
	}
	return s.membersRepo.Delete(ctx, communityID, userID)
}
