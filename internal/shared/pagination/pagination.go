package pagination

const (
	DefaultLimit = 25
	MaxLimit     = 100
)

type Cursor struct {
	After int64
	Limit int
}
