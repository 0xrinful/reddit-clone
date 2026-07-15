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
	Member   MemberDTO `json:"user"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

type MembershipRestrictedDTO struct {
	MembershipDTO
	Permissions int64 `json:"permissions"`
}

// mapping helpers
func toMembershipDTO(m *MembershipView) MembershipDTO {
	return MembershipDTO{
		Member: MemberDTO{
			UserID:    m.UserID,
			Username:  m.Username,
			AvatarUrl: m.AvatarUrl,
		},
		Role:     string(m.Role),
		JoinedAt: m.JoinedAt,
	}
}

func toMembershipRestrictedDTO(m *MembershipView) MembershipRestrictedDTO {
	return MembershipRestrictedDTO{
		MembershipDTO: toMembershipDTO(m),
		Permissions:   int64(m.Permissions),
	}
}

// response envelope
type ListMembersResponse struct {
	Members    any    `json:"members"`
	NextCursor string `json:"next_cursor,omitempty"`
}

type ListModeratorsResponse struct {
	Moderators any    `json:"moderators"`
	NextCursor string `json:"next_cursor,omitempty"`
}

// response constructor
func toList[T any](m []*MembershipView, mapper func(m *MembershipView) T) []T {
	members := make([]T, len(m))
	for i := range m {
		members[i] = mapper(m[i])
	}
	return members
}

func toMembershipList(m []*MembershipView, canViewPermissions bool) any {
	if canViewPermissions {
		return toList(m, toMembershipRestrictedDTO)
	}
	return toList(m, toMembershipDTO)
}

func toListMembersResponse(result *ListResult, nextCursor string) ListMembersResponse {
	return ListMembersResponse{
		Members:    toMembershipList(result.Memberships, result.CanViewPermissions),
		NextCursor: nextCursor,
	}
}

func toListModeratorsResponse(result *ListResult, nextCursor string) ListModeratorsResponse {
	return ListModeratorsResponse{
		Moderators: toMembershipList(result.Memberships, result.CanViewPermissions),
		NextCursor: nextCursor,
	}
}
