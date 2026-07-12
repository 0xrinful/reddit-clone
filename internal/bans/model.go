package bans

import "time"

type BanDuration int

const (
	BanPermanent BanDuration = iota + 1
	BanDay
	BanThreeDays
	BanWeek
	BanMonth
)

func (b BanDuration) Expiry() *time.Time {
	var expiry time.Time
	now := time.Now()

	switch b {
	case BanDay:
		expiry = now.AddDate(0, 0, 1)
	case BanThreeDays:
		expiry = now.AddDate(0, 0, 3)
	case BanWeek:
		expiry = now.AddDate(0, 0, 7)
	case BanMonth:
		expiry = now.AddDate(0, 1, 0)
	case BanPermanent:
		return nil
	default:
		return nil
	}

	return &expiry
}

func (d BanDuration) Valid() bool {
	switch d {
	case BanPermanent, BanDay, BanThreeDays, BanWeek, BanMonth:
		return true
	}
	return false
}

type BanRecord struct {
	CommunityID int64
	UserID      int64
	BannedBy    *int64
	Reason      string
	CreatedAt   time.Time
	ExpiresAt   *time.Time
}

type CreateParams struct {
	CommunityID int64
	Username    string
	Reason      string
	Duration    BanDuration
}

type DeleteParams struct {
	CommunityID int64
	Username    string
}
