package authorization

import (
	"context"

	"github.com/0xrinful/reddit-clone/internal/domain"
)

type Service interface {
	CanBan(ctx context.Context, communityID, actorID, targetID int64) error
	CanUnban(ctx context.Context, communityID, actorID, targetID int64) error
	CanViewBans(ctx context.Context, communityID, actorID int64) error

	CanPost(ctx context.Context, communityID, userID int64) error
	CanViewPost(ctx context.Context, communityID int64, actorID *int64, authorID int64,
		status domain.PostStatus) error
	CanUpdatePost(ctx context.Context, communityID, actorID, authorID int64) error
	CanDeletePost(actorID, authorID int64) error

	CanUpdateCommunity(ctx context.Context, communityID, actorID int64) error
	CanDeleteCommunity(ctx context.Context, communityID, actorID int64) error

	CanManageModerators(ctx context.Context, communityID, actorID, targetID int64) error
	CanViewPermissions(ctx context.Context, communityID, actorID int64) (bool, error)

	CanVotePost(ctx context.Context, actorID, postID int64) error
	CanVoteComment(ctx context.Context, actorID, commentID int64) error

	CanComment(ctx context.Context, actorID, postID int64) error
	CanReply(ctx context.Context, actorID, postID int64, status domain.CommentStatus) error
}

type MemberChecker interface {
	GetAuthority(ctx context.Context, communityID, userID int64) (domain.Authority, error)
	IsMember(ctx context.Context, communityID, userID int64) (bool, error)
}

type BanChecker interface {
	IsBanned(ctx context.Context, communityID, userID int64) (bool, error)
}

type PostsRepo interface {
	GetAuthorizationData(ctx context.Context, id int64) (*domain.PostAuthorizationData, error)
}

type CommentsRepo interface {
	GetAuthorizationData(ctx context.Context, id int64) (*domain.CommentAuthorizationData, error)
}

func NewService(
	banChecker BanChecker,
	memberChecker MemberChecker,
	postsRepo PostsRepo,
	commentsRepo CommentsRepo,
) Service {
	return &service{
		bans:     banChecker,
		members:  memberChecker,
		posts:    postsRepo,
		comments: commentsRepo,
	}
}

type service struct {
	bans     BanChecker
	members  MemberChecker
	posts    PostsRepo
	comments CommentsRepo
}
