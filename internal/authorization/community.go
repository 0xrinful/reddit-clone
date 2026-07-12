package authorization

import (
	"context"

	"github.com/0xrinful/reddit-clone/internal/domain"
	"github.com/0xrinful/reddit-clone/internal/shared/errs"
)

func (s *service) CanUpdateCommunity(ctx context.Context, communityID, actorID int64) error {
	authority, err := s.members.GetAuthority(ctx, communityID, actorID)
	if err != nil {
		return err
	}

	if !authority.Can(domain.PermManageCommunity) {
		return errs.ErrPermissionDenied
	}

	return nil
}

func (s *service) CanDeleteCommunity(ctx context.Context, communityID, actorID int64) error {
	authority, err := s.members.GetAuthority(ctx, communityID, actorID)
	if err != nil {
		return err
	}

	if !authority.Role.IsOwner() {
		return errs.ErrPermissionDenied
	}

	return nil
}
