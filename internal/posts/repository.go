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
	return nil
}

func (r *postgresRepository) Delete(ctx context.Context, id int64) error {
	return nil
}
