package authorization

import (
	"context"

	"github.com/0xrinful/reddit-clone/internal/shared/errs"
)

func (s *service) CanVotePost(ctx context.Context, actorID, postID int64) error {
	info, err := s.posts.GetAuthorizationData(ctx, postID)
	if err != nil {
		return err
	}

	if !info.Status.CanView() {
		return errs.ErrPermissionDenied
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

func (s *service) CanVoteComment(ctx context.Context, actorID, commentID int64) error {
	info, err := s.comments.GetAuthorizationData(ctx, commentID)
	if err != nil {
		return err
	}

	if !info.Status.CanView() {
		return errs.ErrPermissionDenied
	}

	return nil
}
