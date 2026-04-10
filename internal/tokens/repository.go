package tokens

import (
	"context"
	"time"

	"github.com/0xrinful/reddit-clone/internal/database"
)

const (
	ScopeActivation    = "activation"
	ScopeAuth          = "auth"
	ScopePasswordReset = "password-reset"
)

type Repository interface {
	Insert(ctx context.Context, token *Token) error
	DeleteAllForUser(ctx context.Context, scope string, userID int64) error
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
