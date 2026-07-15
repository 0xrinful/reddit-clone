package members

import (
	"time"

	"github.com/0xrinful/reddit-clone/internal/domain"
	"github.com/0xrinful/reddit-clone/internal/shared/validator"
)

// request structs
type PromoteRequest struct {
	Username    string `json:"username"`
	Permissions int64  `json:"permissions"`
}

func (r PromoteRequest) Validate(v *validator.Validator) {
	v.Check(r.Username != "", "username", "must not be empty")
	v.Check(domain.Permission(r.Permissions).Valid(), "permissions", "invalid permissions")
}

// DTOs
type MemberDTO struct {
	UserID    int64   `json:"id"`
	Username  string  `json:"username"`
	AvatarUrl *string `json:"avatar_url"`
}

type MembershipDTO struct {
	Member      MemberDTO `json:"user"`
	Role        string    `json:"role"`
	Permissions int64     `json:"permissions"`
	JoinedAt    time.Time `json:"joined_at"`
}

// mapping helpers
func toMembershipDTO(m *MembershipView) MembershipDTO {
	return MembershipDTO{
		Member: MemberDTO{
			UserID:    m.UserID,
			Username:  m.Username,
			AvatarUrl: m.AvatarUrl,
		},
		Role:        string(m.Role),
		Permissions: int64(m.Permissions),
		JoinedAt:    m.JoinedAt,
	}
}

// response envelope
type ListMembersResponse struct {
	Members    []MembershipDTO `json:"members"`
	NextCursor string          `json:"next_cursor,omitempty"`
}

type ListModeratorsResponse struct {
	Moderators []MembershipDTO `json:"moderators"`
	NextCursor string          `json:"next_cursor,omitempty"`
}

// response constructor
func toListMembersResponse(m []*MembershipView, nextCursor string) ListMembersResponse {
	members := make([]MembershipDTO, len(m))
	for i := range m {
		members[i] = toMembershipDTO(m[i])
	}
	return ListMembersResponse{Members: members, NextCursor: nextCursor}
}

func toListModeratorsResponse(m []*MembershipView, nextCursor string) ListModeratorsResponse {
	moderators := make([]MembershipDTO, len(m))
	for i := range m {
		moderators[i] = toMembershipDTO(m[i])
	}
	return ListModeratorsResponse{Moderators: moderators, NextCursor: nextCursor}
}
