package communities

import (
	"time"

	"github.com/0xrinful/reddit-clone/internal/shared/pagination"
)

type Community struct {
	ID          int64
	Name        string
	OwnerID     *int64
	Description string
	CreatedAt   time.Time
}

type CommunityOwner struct {
	Username string
}

type CommunityView struct {
	Community
	Owner *CommunityOwner
}

type CommunitySummary struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   time.Time
}

type CreateParams struct {
	Name        string
	OwnerID     int64
	Description string
}

type UpdateParams struct {
	Name        *string
	Description *string
}

type ListParams struct {
	Pagination pagination.CursorParams[pagination.CommunityCursor]
}

type SearchParams struct {
	Name       string
	Pagination pagination.OffsetParams
}
