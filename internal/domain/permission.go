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
