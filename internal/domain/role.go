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

func (r Role) IsUser() bool {
	return r == RoleMember
}

func (r Role) CanManageBan(target Role) bool {
	switch r {
	case RoleOwner:
		return target != RoleOwner
	case RoleModerator:
		return target == RoleMember
	default:
		return false
	}
}
