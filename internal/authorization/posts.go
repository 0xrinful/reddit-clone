package authorization

import (
	"context"

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

func (s *service) CanDeletePost(actorID, authorID int64) bool {
	return actorID == authorID
}

func (s *service) CanUpdatePost(actorID, authorID int64) bool {
	return actorID == authorID
}
