package members

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/0xrinful/reddit-clone/internal/database"
	"github.com/0xrinful/reddit-clone/internal/domain"
	"github.com/0xrinful/reddit-clone/internal/shared/errs"
	"github.com/0xrinful/reddit-clone/internal/shared/query"
)

type Repository interface {
	Get(ctx context.Context, communityID, userID int64) (*Membership, error)
	IsMember(ctx context.Context, communityID, userID int64) (bool, error)
	GetAuthority(ctx context.Context, communityID, userID int64) (domain.Authority, error)
	UpdateAuthority(ctx context.Context, communityID, userID int64,
		authority domain.Authority) error
	Create(ctx context.Context, communityID, userID int64, role domain.Role) error
	Delete(ctx context.Context, communityID, userID int64) error
	List(ctx context.Context, p ListParams) ([]*MembershipView, error)
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
		SELECT community_id, user_id, role, permissions, joined_at
		FROM community_members
		WHERE community_id = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := r.db(ctx).QueryRowContext(ctx, query, communityID, userID)
	return scanMembership(row)
}

func (r *postgresRepository) IsMember(
	ctx context.Context,
	communityID, userID int64,
) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM community_members WHERE community_id = $1 AND user_id = $2	
		)`
	var member bool

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.db(ctx).QueryRowContext(ctx, query, communityID, userID).Scan(&member)
	if err != nil {
		return false, err
	}

	return member, nil
}

func (r *postgresRepository) GetAuthority(
	ctx context.Context,
	communityID, userID int64,
) (domain.Authority, error) {
	query := `
		SELECT role, permissions FROM community_members
		WHERE community_id = $1 AND user_id = $2`

	var authority domain.Authority

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.db(ctx).
		QueryRowContext(ctx, query, communityID, userID).
		Scan(&authority.Role, &authority.Permission)
	if err != nil {
		return domain.Authority{}, scanError(err)
	}

	return authority, nil
}

func (r *postgresRepository) Create(
	ctx context.Context,
	communityID, userID int64,
	role domain.Role,
) error {
	query := `
		INSERT INTO community_members (community_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (community_id, user_id) DO NOTHING`

	args := []any{communityID, userID, role}

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

func (r *postgresRepository) List(ctx context.Context, p ListParams) ([]*MembershipView, error) {
	q := query.New()
	cursor := p.Pagination.Cursor

	q.Select(
		"m.community_id, m.user_id, m.role, m.permissions, m.joined_at, u.username, u.avatar_url",
		"community_members m",
	)
	q.Join("users u on m.user_id = u.id")
	q.Where("community_id = ?", p.CommunityID)

	if p.Role != nil {
		q.Where("m.role = ?", p.Role)
	}

	if cursor != nil {
		q.Where("(m.joined_at, m.user_id) < (?, ?)", cursor.JoinedAt, cursor.UserID)
	}
	q.Order("m.joined_at DESC, m.user_id DESC")
	q.Limit(p.Pagination.Limit)

	query, args := q.ToSql()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := r.db(ctx).QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	memberships := make([]*MembershipView, 0, p.Pagination.Limit)
	for rows.Next() {
		m, err := scanView(rows)
		if err != nil {
			return nil, err
		}
		memberships = append(memberships, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return memberships, nil
}

func (r *postgresRepository) UpdateAuthority(
	ctx context.Context,
	communityID, userID int64,
	authority domain.Authority,
) error {
	query := `
		UPDATE community_members SET role = $3, permissions = $4
		WHERE community_id = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	args := []any{communityID, userID, authority.Role, authority.Permission}

	_, err := r.db(ctx).ExecContext(ctx, query, args...)
	return err
}

func scanMembership(row *sql.Row) (*Membership, error) {
	var m Membership

	err := row.Scan(&m.CommunityID, &m.UserID, &m.Role, &m.Permissions, &m.JoinedAt)
	if err != nil {
		return nil, scanError(err)
	}

	return &m, nil
}

func scanView(rows *sql.Rows) (*MembershipView, error) {
	var m MembershipView
	var avatarUrl sql.NullString

	err := rows.Scan(
		&m.CommunityID, &m.UserID,
		&m.Role, &m.Permissions,
		&m.JoinedAt, &m.Username,
		&avatarUrl,
	)
	if err != nil {
		return nil, err
	}

	if avatarUrl.Valid {
		m.AvatarUrl = &avatarUrl.String
	}

	m.JoinedAt = m.JoinedAt.UTC()
	return &m, nil
}

func scanError(err error) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return errs.ErrNotFound
	}
	return err
}
