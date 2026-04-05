package users

import (
	"fmt"
	"time"

	"github.com/0xrinful/reddit-clone/internal/shared/validator"
)

const (
	userNameMin = 3
	userNameMax = 32
	passwordMin = 8
	passwordMax = 72
)

// validation helpers
func ValidateEmail(v *validator.Validator, email string) {
	v.Check(validator.NotBlank(email), "email", "must not be blank")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidateUsername(v *validator.Validator, username string) {
	v.Check(validator.NotBlank(username), "username", "must not be blank")
	v.Check(
		validator.MinLength(username, userNameMin),
		"username",
		fmt.Sprintf("must be at least %d characters", userNameMin),
	)
	v.Check(
		validator.MaxLength(username, userNameMax),
		"username",
		fmt.Sprintf("must not exceed %d characters", userNameMax),
	)
}

func ValidatePassword(v *validator.Validator, password string) {
	v.Check(validator.NotBlank(password), "password", "must not be blank")
	v.Check(
		validator.MinLength(password, passwordMin),
		"password",
		fmt.Sprintf("must be at least %d characters", passwordMin),
	)
	v.Check(
		validator.MaxLength(password, passwordMax),
		"password",
		fmt.Sprintf("must not exceed %d characters", passwordMax),
	)
}

// request structs
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *CreateUserRequest) Validate(v *validator.Validator) {
	ValidateEmail(v, r.Email)
	ValidateUsername(v, r.Username)
	ValidatePassword(v, r.Password)
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
