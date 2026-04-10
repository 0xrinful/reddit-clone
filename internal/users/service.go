package users

type Service interface{}

func NewService(userRepo Repository) Service {
	return &service{userRepo: userRepo}
}

type service struct {
	userRepo Repository
}
