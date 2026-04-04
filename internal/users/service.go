package users

import "context"

type Service interface {
	CreateUser(ctx context.Context, params CreateUserParams) (*User, error)
}

func NewService(repo Repository) Service {
	return &service{repo}
}

type service struct {
	repo Repository
}

func (s *service) CreateUser(ctx context.Context, params CreateUserParams) (*User, error) {
	user := &User{
		Username: params.Username,
		Email:    params.Email,
	}

	err := user.Password.Set(params.PlainPassword)
	if err != nil {
		return nil, err
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
