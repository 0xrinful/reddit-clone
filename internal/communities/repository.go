package communities

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/0xrinful/reddit-clone/internal/shared/errs"
)

type Repository interface {
	GetByName(ctx context.Context, name string) (*Community, error)
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

type postgresRepository struct {
	db *sql.DB
}

func (r *postgresRepository) GetByName(ctx context.Context, name string) (*Community, error) {
	query := `
		SELECT id, name, owner_id, description, created_at, version
		FROM communities 
		WHERE name = $1`

	var c Community
	var ownerID sql.NullInt64

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&c.ID, &c.Name, &ownerID,
		&c.Description, &c.CreatedAt, &c.Version,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	if ownerID.Valid {
		c.OwnerID = &ownerID.Int64
	}
	c.CreatedAt = c.CreatedAt.UTC()

	return &c, nil
}
