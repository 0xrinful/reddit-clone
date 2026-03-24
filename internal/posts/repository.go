package posts

import (
	"context"
	"database/sql"
)

type Repository interface {
	GetByID(ctx context.Context, id int64) (*Post, error)
	Create(ctx context.Context, p *Post) error
	Delete(ctx context.Context, id int64) error
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{db}
}

type postgresRepository struct {
	db *sql.DB
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (*Post, error) {
	return nil, nil
}

func (r *postgresRepository) Create(ctx context.Context, p *Post) error {
	query := `
		INSERT INTO posts (title, body, user_id, community_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []any{p.Title, p.Body, p.UserID, p.CommunityID}

	return r.db.QueryRowContext(ctx, query, args...).Scan(&p.ID, &p.CreatedAt, &p.Version)
}

func (r *postgresRepository) Delete(ctx context.Context, id int64) error {
	return nil
}
