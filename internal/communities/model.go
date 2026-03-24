package communities

import "time"

type Community struct {
	ID          int64
	Name        string
	OwnerID     *int64
	Description string
	CreatedAt   time.Time
	Version     int32
}
