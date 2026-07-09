package bans

import (
	"time"

	"github.com/0xrinful/reddit-clone/internal/shared/validator"
)

type Duration int

type CreateRequest struct {
	Username string `json:"username"`
	Reason   string `json:"reason"`
	Duration int    `json:"duration"`
}

func (c CreateRequest) Validate(v *validator.Validator) {
	if c.Duration > 4 || c.Duration < 1 {
		v.AddError("duration", "invalid ban duration")
	}
	if len(c.Reason) == 0 {
		v.AddError("reason", "must not be empty")
	}
}

func (d Duration) Expiry() *time.Time {
	expiry := time.Now()
	switch d {
	case 1:
		expiry = expiry.Add(BanDay)
	case 2:
		expiry = expiry.Add(BanWeek)
	case 3:
		expiry = expiry.Add(BanMonth)
	default:
		return nil
	}
	return &expiry
}

type DeleteRequest struct {
	Username string `json:"username"`
}
