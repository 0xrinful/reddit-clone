package bans

import (
	"context"
	"time"

	"github.com/0xrinful/reddit-clone/internal/database"
)

type Repository interface {
	IsBanned(ctx context.Context, communityID, userID int64) (bool, error)
	Create(ctx context.Context, b *BanRecord) error
	Delete(ctx context.Context, communityID, userID int64) error
}

func NewRepository(db database.DB) Repository {
	return &postgresRepository{base: db}
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

func (r *postgresRepository) IsBanned(
	ctx context.Context,
	communityID, userID int64,
) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM bans WHERE community_id = $1 AND user_id = $2
		)`

	var banned bool

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.db(ctx).QueryRowContext(ctx, query, communityID, userID).Scan(&banned)
	if err != nil {
		return false, err
	}

	return banned, nil
}

func (r *postgresRepository) Create(ctx context.Context, b *BanRecord) error {
	query := `
		INSERT INTO bans (community_id, user_id, banned_by, reason, expires_at)
		VALUES($1, $2, $3, $4, $5)
		ON CONFLICT (community_id, user_id) DO NOTHING`

	args := []any{b.CommunityID, b.UserID, b.BannedBy, b.Reason, b.ExpiresAt}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := r.db(ctx).ExecContext(ctx, query, args...)
	return err
}

func (r *postgresRepository) Delete(ctx context.Context, communityID, userID int64) error {
	query := `DELETE FROM bans WHERE community_id = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := r.db(ctx).ExecContext(ctx, query, communityID, userID)
	return err
}
