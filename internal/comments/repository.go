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
	GetByID(ctx context.Context, id int64) (*Comment, error)
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

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (*Comment, error) {
	query := `
		SELECT
			id, post_id, user_id, parent_id, content, score, reply_count, status, created_at
		FROM comments
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := r.db(ctx).QueryRowContext(ctx, query, id)

	return scanComment(row)
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

func scanComment(row *sql.Row) (*Comment, error) {
	var c Comment
	var parentID sql.NullInt64

	err := row.Scan(
		&c.ID, &c.PostID, &c.UserID, &parentID,
		&c.Content, &c.Score, &c.ReplyCount, &c.Status, &c.CreatedAt,
	)
	if err != nil {
		return nil, scanError(err)
	}

	if parentID.Valid {
		c.ParentID = &parentID.Int64
	}
	c.CreatedAt = c.CreatedAt.UTC()
	return &c, nil
}

func scanError(err error) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return errs.ErrNotFound
	}
	return err
}
