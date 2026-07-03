package members

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/0xrinful/reddit-clone/internal/database"
	"github.com/0xrinful/reddit-clone/internal/shared/errs"
)

type Repository interface {
	Get(ctx context.Context, communityID, userID int64) (*Membership, error)
	Create(ctx context.Context, m *Membership) error
	Delete(ctx context.Context, communityID, userID int64) error
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

func (r *postgresRepository) Get(
	ctx context.Context,
	communityID, userID int64,
) (*Membership, error) {
	query := `
		SELECT community_id, user_id, role, joined_at
		FROM community_members
		WHERE community_id = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := r.db(ctx).QueryRowContext(ctx, query, communityID, userID)
	return scanMembership(row)
}

func (r *postgresRepository) Create(ctx context.Context, m *Membership) error {
	query := `
		INSERT INTO community_members (community_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (community_id, user_id) DO NOTHING`

	args := []any{m.CommunityID, m.UserID, m.Role}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := r.db(ctx).ExecContext(ctx, query, args...)
	return err
}

func (r *postgresRepository) Delete(ctx context.Context, communityID, userID int64) error {
	query := `DELETE FROM community_members WHERE community_id = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := r.db(ctx).ExecContext(ctx, query, communityID, userID)
	return err
}

func scanMembership(row *sql.Row) (*Membership, error) {
	var m Membership

	err := row.Scan(&m.CommunityID, &m.UserID, &m.Role, &m.JoinedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &m, nil
}
