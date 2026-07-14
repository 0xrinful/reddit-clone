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

func (r Role) Level() int {
	switch r {
	case RoleMember:
		return 1
	case RoleModerator:
		return 2
	case RoleOwner:
		return 3
	default:
		return 0
	}
}

func (r Role) AtLeast(role Role) bool {
	return r.Level() >= role.Level()
}
