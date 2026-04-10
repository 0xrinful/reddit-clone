package users

import (
	"time"
)

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
func ToUserPublicDTO(u *User) UserPublicDTO {
	return UserPublicDTO{
		Username:  u.Username,
		CreatedAt: u.CreatedAt,
		AvatarUrl: u.AvatarUrl,
	}
}

func ToUserOwnerDTO(u *User) UserOwnerDTO {
	return UserOwnerDTO{
		UserPublicDTO: ToUserPublicDTO(u),
		ID:            u.ID,
		Email:         u.Email,
		Activated:     u.Activated,
		ActivatedAt:   u.ActivatedAt,
	}
}
