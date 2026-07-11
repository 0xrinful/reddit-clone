package authorization

import "github.com/0xrinful/reddit-clone/internal/domain"

func (s *service) CanBan(actor, target domain.Authority) bool {
	if !actor.Can(domain.PermBanUsers) {
		return false
	}

	if target.Role.IsOwner() {
		return false
	}

	if target.Role.IsModerator() && !actor.Can(domain.PermManageModerators) {
		return false
	}

	return true
}

func (s *service) CanUnban(actor, target domain.Authority) bool {
	return s.CanBan(actor, target)
}
