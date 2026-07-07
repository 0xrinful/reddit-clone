package members

import "time"

// request structs

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

// response envelope
type ListMembersResponse struct {
	Members    []MembershipDTO `json:"members"`
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
