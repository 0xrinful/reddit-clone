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
