package authorization

import (
	"context"

	"github.com/0xrinful/reddit-clone/internal/domain"
	"github.com/0xrinful/reddit-clone/internal/shared/errs"
)

func (s *service) CanBan(ctx context.Context, communityID, actorID, targetID int64) error {
	actor, err := s.members.GetAuthority(ctx, communityID, actorID)
	if err != nil {
		return err
	}

	if !actor.Can(domain.PermBanUsers) {
		return errs.ErrPermissionDenied
	}

	if actorID == targetID {
		return errs.ErrSelfBan
	}

	target, err := s.members.GetAuthority(ctx, communityID, targetID)
	if err != nil {
		return err
	}

	if target.Role.IsOwner() {
		return errs.ErrPermissionDenied
	}

	if target.Role.IsModerator() && !actor.Can(domain.PermManageModerators) {
		return errs.ErrPermissionDenied
	}

	return nil
}

func (s *service) CanUnban(ctx context.Context, communityID, actorID, targetID int64) error {
	actor, err := s.members.GetAuthority(ctx, communityID, actorID)
	if err != nil {
		return err
	}

	if !actor.Can(domain.PermBanUsers) {
		return errs.ErrPermissionDenied
	}

	target, err := s.members.GetAuthority(ctx, communityID, targetID)
	if err != nil {
		return err
	}

	if target.Role.IsOwner() {
		return errs.ErrPermissionDenied
	}

	if target.Role.IsModerator() && !actor.Can(domain.PermManageModerators) {
		return errs.ErrPermissionDenied
	}

	return nil
}
