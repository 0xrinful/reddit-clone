package domain

type Authority struct {
	Role       Role
	Permission Permission
}

func (a Authority) Can(p Permission) bool {
	if a.Role == RoleOwner {
		return true
	}
	return a.Permission.Has(p)
}

// TODO: move this to Permission service later
func (a Authority) CanActOn(target Authority) bool {
	switch a.Role {
	case RoleOwner:
		return !target.Role.IsOwner()
	case RoleModerator:
		return target.Role.IsMember()
	default:
		return false
	}
}
