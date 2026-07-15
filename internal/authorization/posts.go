package authorization

import (
	"context"

	"github.com/0xrinful/reddit-clone/internal/domain"
	"github.com/0xrinful/reddit-clone/internal/shared/errs"
)

func (s *service) CanPost(ctx context.Context, communityID, userID int64) error {
	banned, err := s.bans.IsBanned(ctx, communityID, userID)
	if err != nil {
		return err
	}

	if banned {
		return errs.ErrBanned
	}

	member, err := s.members.IsMember(ctx, communityID, userID)
	if err != nil {
		return err
	}

	if !member {
		return errs.ErrNotMember
	}

	return nil
}

func (s *service) CanViewPost(ctx context.Context, communityID int64, actorID *int64,
	authorID int64, status domain.PostStatus,
) error {
	if status.IsVisible() {
		return nil
	}

	if actorID == nil {
		return errs.ErrPermissionDenied
	}
	id := *actorID

	if id == authorID {
		return nil
	}

	actor, err := s.members.GetAuthority(ctx, communityID, id)
	if err != nil {
		return err
	}

	if actor.Role.AtLeast(domain.RoleModerator) {
		return nil
	}

	return errs.ErrPermissionDenied
}

func (s *service) CanDeletePost(actorID, authorID int64) error {
	if actorID != authorID {
		return errs.ErrPermissionDenied
	}
	return nil
}

func (s *service) CanUpdatePost(ctx context.Context, communityID, actorID, authorID int64) error {
	if actorID != authorID {
		return errs.ErrPermissionDenied
	}

	banned, err := s.bans.IsBanned(ctx, communityID, authorID)
	if err != nil {
		return err
	}

	if banned {
		return errs.ErrBanned
	}

	return nil
}
