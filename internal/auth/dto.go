package auth

import (
	"github.com/0xrinful/reddit-clone/internal/shared/validator"
	"github.com/0xrinful/reddit-clone/internal/users"
)

// request structs
type RegisterUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *RegisterUserRequest) Validate(v *validator.Validator) {
	validator.ValidateEmail(v, r.Email)
	validator.ValidateUsername(v, r.Username)
	validator.ValidatePassword(v, r.Password)
}

type ActivateUserRequest struct {
	Token string `json:"token"`
}

func (r *ActivateUserRequest) Validate(v *validator.Validator) {
	v.Check(r.Token != "", "token", "must be provided")
}

// response envelope
type RegisterUserResponse struct {
	User users.UserOwnerDTO `json:"user"`
}

// response constructor
func toRegisterUserResponse(u *users.User) RegisterUserResponse {
	return RegisterUserResponse{User: users.ToUserOwnerDTO(u)}
}
