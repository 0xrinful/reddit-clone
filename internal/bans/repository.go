package bans

import (
	"context"
	"database/sql"
	"time"

	"github.com/0xrinful/reddit-clone/internal/database"
	"github.com/0xrinful/reddit-clone/internal/shared/query"
)

type Repository interface {
	IsBanned(ctx context.Context, communityID, userID int64) (bool, error)
	Create(ctx context.Context, b *BanRecord) error
	Delete(ctx context.Context, communityID, userID int64) error
	List(ctx context.Context, p ListParams) ([]*BanView, error)
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
			SELECT 1 FROM community_bans 
			WHERE community_id = $1 AND user_id = $2 AND (expires_at IS NULL OR expires_at > now())
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
		INSERT INTO community_bans (community_id, user_id, banned_by, reason, expires_at)
		VALUES($1, $2, $3, $4, $5)
		ON CONFLICT (community_id, user_id) DO NOTHING`

	args := []any{b.CommunityID, b.UserID, b.BannedBy, b.Reason, b.ExpiresAt}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := r.db(ctx).ExecContext(ctx, query, args...)
	return err
}

func (r *postgresRepository) Delete(ctx context.Context, communityID, userID int64) error {
	query := `DELETE FROM community_bans WHERE community_id = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := r.db(ctx).ExecContext(ctx, query, communityID, userID)
	return err
}

func (r *postgresRepository) List(ctx context.Context, p ListParams) ([]*BanView, error) {
	q := query.New()
	cursor := p.Pagination.Cursor

	q.Select(
		`b.community_id, b.user_id, b.banned_by, b.reason, b.created_at, 
	  b.expires_at, u.username, u.avatar_url, bu.username, bu.avatar_url`,
		"community_bans b",
	)
	q.Join("users u on b.user_id = u.id")
	q.LeftJoin("users bu on b.banned_by = bu.id")
	q.Where("b.community_id = ?", p.CommunityID)

	if cursor != nil {
		q.Where("(b.created_at, b.user_id) < (?, ?)", cursor.CreatedAt, cursor.UserID)
	}
	q.Order("b.created_at DESC, b.user_id DESC")
	q.Limit(p.Pagination.Limit)

	query, args := q.ToSql()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := r.db(ctx).QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bans := make([]*BanView, 0, p.Pagination.Limit)
	for rows.Next() {
		b, err := scanView(rows)
		if err != nil {
			return nil, err
		}
		bans = append(bans, b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bans, nil
}

func scanView(rows *sql.Rows) (*BanView, error) {
	var b BanView
	var userAvatarUrl sql.NullString
	var expiresAt sql.NullTime
	var modID sql.NullInt64
	var modUsername sql.NullString
	var modAvatarUrl sql.NullString

	err := rows.Scan(
		&b.CommunityID, &b.UserID, &modID,
		&b.Reason, &b.CreatedAt, &expiresAt,
		&b.BannedUser.Username, &userAvatarUrl, &modUsername, &modAvatarUrl,
	)
	if err != nil {
		return nil, err
	}
	b.CreatedAt = b.CreatedAt.UTC()

	if userAvatarUrl.Valid {
		b.BannedUser.AvatarUrl = &userAvatarUrl.String
	}

	if expiresAt.Valid {
		t := expiresAt.Time.UTC()
		b.ExpiresAt = &t
	}

	if modID.Valid {
		b.BannedBy = &modID.Int64
		b.Moderator = &UserView{}
		b.Moderator.Username = modUsername.String

		if modAvatarUrl.Valid {
			b.Moderator.AvatarUrl = &modAvatarUrl.String
		}
	}

	return &b, nil
}
