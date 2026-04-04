package users

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"

	"github.com/0xrinful/reddit-clone/internal/shared/errs"
)

type Repository interface {
	Create(ctx context.Context, u *User) error
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

type postgresRepository struct {
	db *sql.DB
}

func (r *postgresRepository) Create(ctx context.Context, u *User) error {
	query := `
		INSERT INTO users (username, email, hashed_password) 
		VALUES ($1, $2, $3)
		RETURNING created_at, version, activated`

	args := []any{u.ID, u.Username, u.Email, u.Password.hash}

	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&u.CreatedAt, &u.Version, u.Activated)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "users_email_key":
				return errs.ErrDuplicateEmail
			case "users_username_key":
				return errs.ErrDuplicateUsername
			}
		}
		return err
	}
	return nil
}
