package posts

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/0xrinful/reddit-clone/internal/shared/apperr"
)

type Repository interface {
	Get(ctx context.Context, id int64) (*Post, error)
	Create(ctx context.Context, p *Post) error
	Delete(ctx context.Context, id int64) error
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{db}
}

type postgresRepository struct {
	db *sql.DB
}

func (r *postgresRepository) Get(ctx context.Context, id int64) (*Post, error) {
	query := `
		SELECT id, title, body, user_id, community_id, views, created_at, version
		FROM posts
		WHERE id = $1`

	var p Post

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Title, &p.Body, &p.UserID, &p.CommunityID,
		&p.Views, &p.CreatedAt, &p.Version,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *postgresRepository) Create(ctx context.Context, p *Post) error {
	query := `
		INSERT INTO posts (title, body, user_id, community_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []any{p.Title, p.Body, p.UserID, p.CommunityID}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return r.db.QueryRowContext(ctx, query, args...).Scan(&p.ID, &p.CreatedAt, &p.Version)
}

func (r *postgresRepository) Delete(ctx context.Context, id int64) error {
	return nil
}
