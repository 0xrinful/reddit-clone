package communities

import (
	"time"

	"github.com/0xrinful/reddit-clone/internal/shared/validator"
)

// request structs
type CreateCommunitytRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r *CreateCommunitytRequest) Validate(v *validator.Validator) {
	validator.ValidateCommunityName(v, r.Name)
	validator.ValidateCommunityDescription(v, r.Description)
}

type UpdateCommunitytRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func (r *UpdateCommunitytRequest) Validate(v *validator.Validator) {
	if r.Name == nil && r.Description == nil {
		v.AddError("request", "must provide at least one field")
		return
	}
	if r.Name != nil {
		validator.ValidateCommunityName(v, *r.Name)
	}
	if r.Description != nil {
		validator.ValidateCommunityDescription(v, *r.Description)
	}
}

// DTOs
type OwnerDTO struct {
	ID       int64  `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}

type CommunityDTO struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Owner       *OwnerDTO `json:"owner"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type CommunitySummaryDTO struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// mapping helpers
func communityToView(c *Community) *CommunityView {
	return &CommunityView{
		Community: *c,
		Owner:     &CommunityOwner{},
	}
}

func toCommunityDTO(c *CommunityView) CommunityDTO {
	dto := CommunityDTO{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
	}

	if c.OwnerID != nil {
		dto.Owner = &OwnerDTO{
			ID:       *c.OwnerID,
			Username: c.Owner.Username,
		}
	}
	return dto
}

func toCommunitySummaryDTO(c *CommunitySummary) CommunitySummaryDTO {
	return CommunitySummaryDTO{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
	}
}

// response envelope
type CommunityResponse struct {
	Community CommunityDTO `json:"community"`
}

type SearchCommunitiesResponse struct {
	Communities []CommunitySummaryDTO `json:"communities"`
}

// response constructor
func toCommunityResponse(c *CommunityView) CommunityResponse {
	return CommunityResponse{
		Community: toCommunityDTO(c),
	}
}

func toSearchCommunitiesResponse(c []*CommunitySummary) SearchCommunitiesResponse {
	communities := make([]CommunitySummaryDTO, len(c))
	for i := range c {
		communities[i] = toCommunitySummaryDTO(c[i])
	}
	return SearchCommunitiesResponse{Communities: communities}
}
