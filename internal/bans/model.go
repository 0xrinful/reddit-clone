package bans

import "time"

var (
	BanDay                      = 24 * time.Hour
	BanWeek                     = 7 * BanDay
	BanMonth                    = 30 * BanDay
	BanPermanent *time.Duration = nil
)

type BanRecord struct {
	UserID      int64
	CommunityID int64
	BannedBy    int64
	BannedAt    *time.Time
	Expiry      *time.Time
	Reason      string
}

type CreateParams struct {
	UserID      int64
	CommunityID int64
	Username    string
	Reason      string
	Duration    Duration
}

type DeleteParams struct {
	UserID      int64
	CommunityID int64
	Username    string
}
