package domain

type Permission int64

const (
	PermBanUsers Permission = 1 << iota
	PermManageModerators
	PermDeletePosts
)

func (p Permission) Has(target Permission) bool {
	return p&target != 0
}
