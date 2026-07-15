package members

import (
	"time"

	"github.com/0xrinful/reddit-clone/internal/domain"
	"github.com/0xrinful/reddit-clone/internal/shared/pagination"
)

type Member struct {
	Username  string
	AvatarUrl *string
}

type Membership struct {
	UserID      int64
	CommunityID int64
	Role        domain.Role
	Permissions domain.Permission
	JoinedAt    time.Time
}

type MembershipView struct {
	Membership
	Member
}

type ListParams struct {
	CommunityID int64
	Role        *domain.Role
	Pagination  pagination.CursorParams[pagination.MemberCursor]
}

type ListResult struct {
	Memberships        []*MembershipView
	CanViewPermissions bool
}

type PromoteParams struct {
	CommunityID int64
	Username    string
	Perms       domain.Permission
}

type DemoteParams struct {
	CommunityID int64
	Username    string
}
