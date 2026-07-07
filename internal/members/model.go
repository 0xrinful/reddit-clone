package members

import (
	"time"

	"github.com/0xrinful/reddit-clone/internal/shared/pagination"
)

type CommunityRole string

const (
	RoleMember    CommunityRole = "member"
	RoleModerator CommunityRole = "moderator"
	RoleOwner     CommunityRole = "owner"
)

type Member struct {
	Username  string
	AvatarUrl *string
}

type Membership struct {
	UserID      int64
	CommunityID int64
	Role        CommunityRole
	JoinedAt    time.Time
}

type MembershipView struct {
	Membership
	Member
}

type ListParams struct {
	CommunityID int64
	Pagination  pagination.CursorParams[pagination.MemberCursor]
}
