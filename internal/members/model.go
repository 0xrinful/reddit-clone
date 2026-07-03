package members

import "time"

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
