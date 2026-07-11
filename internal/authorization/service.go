package authorization

import (
	"context"

	"github.com/0xrinful/reddit-clone/internal/domain"
)

type Service interface {
	CanBan(actor, target domain.Authority) bool
	CanUnban(actor, target domain.Authority) bool

	CanPost(ctx context.Context, communityID, userID int64) error
	CanDeletePost(actorID, authorID int64) bool
	CanModifyPost(actorID, authorID int64) bool
}

type MemberChecker interface {
	GetAuthority(ctx context.Context, communityID, userID int64) (domain.Authority, error)
	IsMember(ctx context.Context, communityID, userID int64) (bool, error)
}

type BanChecker interface {
	IsBanned(ctx context.Context, communityID, userID int64) (bool, error)
}

func NewService(banChecker BanChecker, memberChecker MemberChecker) Service {
	return &service{bans: banChecker, members: memberChecker}
}

type service struct {
	bans    BanChecker
	members MemberChecker
}
