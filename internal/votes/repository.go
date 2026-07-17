package votes

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/0xrinful/reddit-clone/internal/database"
)

type Repository interface {
	CreateInitialPostVote(ctx context.Context, userID, postID int64) error
	VotePost(ctx context.Context, vote PostVote) (int64, error)
	UnvotePost(ctx context.Context, userID, postID int64) (int64, error)
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

func (r *postgresRepository) CreateInitialPostVote(
	ctx context.Context,
	userID, postID int64,
) error {
	query := `
		INSERT INTO post_votes (user_id, post_id, value)
		VALUES ($1, $2, 1)`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := r.db(ctx).ExecContext(ctx, query, userID, postID)
	return err
}

func (r *postgresRepository) VotePost(ctx context.Context, vote PostVote) (int64, error) {
	query := `
		INSERT INTO post_votes (user_id, post_id, value)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, post_id)
		DO UPDATE SET value = EXCLUDED.value
		RETURNING new.value - COALESCE(old.value, 0) AS delta`

	var delta int64

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	args := []any{vote.UserID, vote.PostID, vote.Value}

	err := r.db(ctx).QueryRowContext(ctx, query, args...).Scan(&delta)
	if err != nil {
		return 0, err
	}

	return delta, nil
}

func (r *postgresRepository) UnvotePost(ctx context.Context, userID, postID int64) (int64, error) {
	query := `
		DELETE FROM post_votes 
		WHERE user_id = $1 AND post_id = $2
		RETURNING -value AS delta`

	var delta int64

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.db(ctx).QueryRowContext(ctx, query, userID, postID).Scan(&delta)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return delta, nil
}
