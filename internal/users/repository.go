package users

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"

	"github.com/0xrinful/reddit-clone/internal/database"
	"github.com/0xrinful/reddit-clone/internal/shared/errs"
)

type Repository interface {
	Create(ctx context.Context, u *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	SetActivated(ctx context.Context, userID int64) error
}

func NewRepository(db database.DB) Repository {
	return &postgresRepository{db: db}
}

type postgresRepository struct {
	db database.DB
}

func (r *postgresRepository) Create(ctx context.Context, u *User) error {
	query := `
		INSERT INTO users (username, email, hashed_password) 
		VALUES ($1, $2, $3)
		RETURNING id, created_at, version, activated`

	args := []any{u.Username, u.Email, u.Password.hash}

	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).
		Scan(&u.ID, &u.CreatedAt, &u.Version, &u.Activated)
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

func (r *postgresRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, username, email, hashed_password, avatar_url, 
		created_at, version, activated, activated_at
		FROM users WHERE email = $1`

	var user User
	var avatarUrl sql.NullString
	var activatedAt sql.NullTime

	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password.hash,
		&avatarUrl, &user.CreatedAt, &user.Version, &user.Activated,
		&activatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	if avatarUrl.Valid {
		user.AvatarUrl = &avatarUrl.String
	}

	if activatedAt.Valid {
		activatedAt := activatedAt.Time.UTC()
		user.ActivatedAt = &activatedAt
	}

	user.CreatedAt = user.CreatedAt.UTC()
	return &user, nil
}

func (r *postgresRepository) SetActivated(ctx context.Context, userID int64) error {
	query := `
		UPDATE users SET activated = true, activated_at = NOW()
		WHERE id = $1 and activated = false`

	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
