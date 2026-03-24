package communities

import (
	"context"
	"database/sql"
	"errors"

	"github.com/0xrinful/reddit-clone/internal/shared/apperr"
)

type Repository interface {
	GetByName(ctx context.Context, name string) (*Community, error)
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{}
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

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&c.ID, &c.Name, &c.OwnerID,
		&c.Description, &c.CreatedAt, &c.Version,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &c, nil
}
