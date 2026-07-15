package bans

import (
	"time"

	"github.com/0xrinful/reddit-clone/internal/shared/validator"
)

// request structs
type CreateRequest struct {
	Username string `json:"username"`
	Reason   string `json:"reason"`
	Duration int    `json:"duration"`
}

func (r *CreateRequest) Validate(v *validator.Validator) {
	v.Check(r.Username != "", "username", "must not be empty")
	v.Check(BanDuration(r.Duration).Valid(), "duration", "invalid ban duration")
	v.Check(r.Reason != "", "reason", "must not be empty")
}

// DTOs
type UserDTO struct {
	UserID    int64   `json:"id"`
	Username  string  `json:"username"`
	AvatarUrl *string `json:"avatar_url"`
}

type BanDTO struct {
	User      UserDTO    `json:"user"`
	Moderator *UserDTO   `json:"moderator"`
	Reason    string     `json:"reason"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at"`
}

// mapping helpers
func toBanDTO(b *BanView) BanDTO {
	dto := BanDTO{
		User: UserDTO{
			UserID:    b.UserID,
			Username:  b.BannedUser.Username,
			AvatarUrl: b.BannedUser.AvatarUrl,
		},
		Reason:    b.Reason,
		CreatedAt: b.CreatedAt,
		ExpiresAt: b.ExpiresAt,
	}

	if b.BannedBy != nil {
		dto.Moderator = &UserDTO{
			UserID:    *b.BannedBy,
			Username:  b.Moderator.Username,
			AvatarUrl: b.Moderator.AvatarUrl,
		}
	}
	return dto
}

// response envelope
type ListBansResponse struct {
	Bans       []BanDTO `json:"bans"`
	NextCursor string   `json:"next_cursor,omitempty"`
}

// response constructor
func toListBansResponse(b []*BanView, nextCursor string) ListBansResponse {
	bans := make([]BanDTO, len(b))
	for i := range b {
		bans[i] = toBanDTO(b[i])
	}
	return ListBansResponse{Bans: bans, NextCursor: nextCursor}
}
