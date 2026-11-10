package authorization

import (
	"context"

	"github.com/0xrinful/reddit-clone/internal/domain"
	"github.com/0xrinful/reddit-clone/internal/shared/errs"
)

func (s *service) CanComment(ctx context.Context, actorID, postID int64) error {
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
		return errs.ErrBanned
	}

	return nil
}

func (s *service) CanReply(ctx context.Context, actorID, postID int64,
	status domain.CommentStatus,
) error {
	if err := s.CanComment(ctx, actorID, postID); err != nil {
		return err
	}

	if !status.CanView() {
		return errs.ErrPermissionDenied
	}

	return nil
}

func (s *service) CanUpdateComment(
	ctx context.Context,
	actorID int64,
	info *domain.CommentAuthorizationData,
) error {
	if actorID != info.AuthorID {
		return errs.ErrPermissionDenied
	}

	banned, err := s.bans.IsBanned(ctx, info.CommunityID, actorID)
	if err != nil {
		return err
	}

	if banned {
		return errs.ErrBanned
	}

	return nil
}

func (s *service) CanDeleteComment(ctx context.Context, actorID, authorID int64) error {
	if actorID != authorID {
		return errs.ErrPermissionDenied
	}

	return nil
}
