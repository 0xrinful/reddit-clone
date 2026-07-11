package bans

import (
	"github.com/0xrinful/reddit-clone/internal/shared/validator"
)

// request structs
type CreateRequest struct {
	Username string `json:"username"`
	Reason   string `json:"reason"`
	Duration int    `json:"duration"`
}

func (c CreateRequest) Validate(v *validator.Validator) {
	v.Check(BanDuration(c.Duration).Valid(), "duration", "invalid ban duration")
	v.Check(c.Reason != "", "reason", "must not be empty")
}

// DTOs

// mapping helpers

// response envelope

// response constructor
