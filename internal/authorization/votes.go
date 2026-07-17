package authorization

import (
	"context"

	"github.com/0xrinful/reddit-clone/internal/shared/errs"
)

func (s *service) CanVotePost(ctx context.Context, actorID, postID int64) error {
	info, err := s.posts.GetAuthorizationInfo(ctx, postID)
	if err != nil {
		return err
	}

	banned, err := s.bans.IsBanned(ctx, info.CommunityID, actorID)
	if err != nil {
		return err
	}

	if banned {
		return errs.ErrPermissionDenied
	}

	return nil
}
