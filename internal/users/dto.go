package users

import (
	"time"

	"github.com/0xrinful/reddit-clone/internal/shared/validator"
)

// request structs
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *CreateUserRequest) Validate(v *validator.Validator) {
	validator.ValidateEmail(v, r.Email)
	validator.ValidateUsername(v, r.Username)
	validator.ValidatePassword(v, r.Password)
}

type ActivateUserRequest struct {
	Token string `json:"token"`
}

func (r *ActivateUserRequest) Validate(v *validator.Validator) {
	validator.ValidateToken(v, r.Token)
}

// DTOs
type UserPublicDTO struct {
	Username  string    `json:"username"`
	AvatarUrl *string   `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
}

type UserOwnerDTO struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	UserPublicDTO
	Activated   bool       `json:"activated"`
	ActivatedAt *time.Time `json:"activated_at,omitempty"`
}

// mapping helpers
func toUserPublicDTO(u *User) UserPublicDTO {
	return UserPublicDTO{
		Username:  u.Username,
		CreatedAt: u.CreatedAt,
		AvatarUrl: u.AvatarUrl,
	}
}

func toUserOwnerDTO(u *User) UserOwnerDTO {
	return UserOwnerDTO{
		UserPublicDTO: toUserPublicDTO(u),
		ID:            u.ID,
		Email:         u.Email,
		Activated:     u.Activated,
		ActivatedAt:   u.ActivatedAt,
	}
}

// response envelope
type UserResponse struct {
	User UserPublicDTO `json:"user"`
}

type UserOwnerResponse struct {
	User UserOwnerDTO `json:"user"`
}

// response constructor
func toUserResponse(u *User) UserResponse {
	return UserResponse{User: toUserPublicDTO(u)}
}

func toOwnerResponse(u *User) UserOwnerResponse {
	return UserOwnerResponse{User: toUserOwnerDTO(u)}
}
