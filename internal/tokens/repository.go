package tokens

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/0xrinful/reddit-clone/internal/database"
	"github.com/0xrinful/reddit-clone/internal/shared/errs"
)

const (
	ScopeActivation    = "activation"
	ScopeAuth          = "auth"
	ScopePasswordReset = "password-reset"
)

type Repository interface {
	Insert(ctx context.Context, token *Token) error
	DeleteAllForUser(ctx context.Context, scope string, userID int64) error
	GetByHash(ctx context.Context, scope string, hashed string) (*Token, error)
}

func NewRepository(db database.DB) Repository {
	return &postgresRepository{db}
}

type postgresRepository struct {
	db database.DB
}

func (r *postgresRepository) Insert(ctx context.Context, token *Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope) 
		VALUES ($1, $2, $3, $4)`

	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *postgresRepository) DeleteAllForUser(
	ctx context.Context,
	scope string,
	userID int64,
) error {
	query := `DELETE FROM tokens WHERE user_id = $1 AND scope = $2`

	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, userID, scope)
	return err
}

func (r *postgresRepository) GetByHash(
	ctx context.Context,
	scope string,
	hashed string,
) (*Token, error) {
	query := `
		SELECT hash, user_id, expiry, scope FROM tokens
		WHERE hash = $1 AND scope = $2`

	var token Token

	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, hashed, scope).
		Scan(&token.Hash, &token.UserID, &token.Expiry, &token.Scope)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrInvalidToken
	}

	if err != nil {
		return nil, err
	}

	token.Expiry = token.Expiry.UTC()
	return &token, nil
}
