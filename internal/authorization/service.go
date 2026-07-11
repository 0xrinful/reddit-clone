package authorization

import "github.com/0xrinful/reddit-clone/internal/domain"

type Service interface {
	CanBan(actor, target domain.Authority) bool
}

func NewService() Service {
	return &service{}
}

type service struct{}
