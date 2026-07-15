package domain

type Permission int64

const (
	PermRemovePosts Permission = 1 << iota
	PermBanUsers
	PermManageModerators
	PermManageCommunity
)

func (p Permission) Has(target Permission) bool {
	return p&target != 0
}

func (p Permission) Valid() bool {
	return p&^(PermRemovePosts|PermBanUsers|PermManageModerators|PermManageCommunity) == 0
}

func DefaultPerms(role Role) Permission {
	switch role {
	case RoleModerator:
		return PermRemovePosts | PermBanUsers
	case RoleMember:
		return 0
	default:
		return 0
	}
}
