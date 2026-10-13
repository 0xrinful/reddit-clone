package comments

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/0xrinful/reddit-clone/internal/database"
	"github.com/0xrinful/reddit-clone/internal/shared/errs"
)

type Repository interface {
	Create(ctx context.Context, c *Comment) error
}

func NewRepository(db database.DB) Repository {
	return &postgresRepository{db}
}

type postgresRepository struct {
	base database.DB
}

func (r *postgresRepository) db(ctx context.Context) database.DB {
	if tx, ok := database.GetTx(ctx); ok {
		return tx
	}
	return r.base
}

func (r *postgresRepository) Create(ctx context.Context, c *Comment) error {
	query := `
		INSERT INTO comments (post_id, user_id, parent_id, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id, depth, score, reply_count, created_at`

	args := []any{c.PostID, c.UserID, c.ParentID, c.Content}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.db(ctx).
		QueryRowContext(ctx, query, args...).
		Scan(&c.ID, &c.Depth, &c.Score, &c.ReplyCount, &c.CreatedAt)
	if err != nil {
		return err
	}
	c.CreatedAt = c.CreatedAt.UTC()

	return nil
}

func scanError(err error) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return errs.ErrNotFound
	}
	return err
}
