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
	v.Check(r.Token != "", "token", "is required")
}

type SendActivationEmailRequest struct {
	Email string `json:"email"`
}

func (r *SendActivationEmailRequest) Validate(v *validator.Validator) {
	validator.ValidateEmail(v, r.Email)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate(v *validator.Validator) {
	validator.ValidateEmail(v, r.Email)
	validator.ValidatePassword(v, r.Password)
}

type RefreshRequest struct {
	Token string `json:"refresh_token"`
}

func (r *RefreshRequest) Validate(v *validator.Validator) {
	v.Check(r.Token != "", "refresh_token", "is required")
}

type LogoutRequest struct {
	Token string `json:"refresh_token"`
}

func (r *LogoutRequest) Validate(v *validator.Validator) {
	v.Check(r.Token != "", "refresh_token", "is required")
}

// response envelope
type RegisterUserResponse struct {
	User users.UserOwnerDTO `json:"user"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// response constructor
func toRegisterUserResponse(u *users.User) RegisterUserResponse {
	return RegisterUserResponse{User: users.ToUserOwnerDTO(u)}
}

func toLoginResponse(access, refresh string) LoginResponse {
	return LoginResponse{AccessToken: access, RefreshToken: refresh}
}

func toRefreshResponse(access, refresh string) RefreshResponse {
	return RefreshResponse{AccessToken: access, RefreshToken: refresh}
}
