package authorization

import (
	"context"

	"github.com/0xrinful/reddit-clone/internal/domain"
	"github.com/0xrinful/reddit-clone/internal/shared/errs"
)

func (s *service) CanManageModerators(
	ctx context.Context,
	communityID, actorID, targetID int64,
) error {
	actor, err := s.members.GetAuthority(ctx, communityID, actorID)
	if err != nil {
		return err
	}

	if !actor.Can(domain.PermManageModerators) {
		return errs.ErrPermissionDenied
	}

	target, err := s.members.GetAuthority(ctx, communityID, targetID)
	if err != nil {
		return err
	}

	if target.Role.IsOwner() {
		return errs.ErrPermissionDenied
	}

	return nil
}
