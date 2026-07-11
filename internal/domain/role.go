package domain

type Role string

const (
	RoleOwner     Role = "owner"
	RoleModerator Role = "moderator"
	RoleMember    Role = "member"
)

func (r Role) IsOwner() bool {
	return r == RoleOwner
}

func (r Role) IsModerator() bool {
	return r == RoleModerator
}

func (r Role) IsMember() bool {
	return r == RoleMember
}
