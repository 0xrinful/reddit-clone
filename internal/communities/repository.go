package communities

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"

	"github.com/0xrinful/reddit-clone/internal/database"
	"github.com/0xrinful/reddit-clone/internal/shared/errs"
	"github.com/0xrinful/reddit-clone/internal/shared/query"
)

type Repository interface {
	GetByName(ctx context.Context, name string) (*Community, error)
	GetViewByName(ctx context.Context, name string) (*CommunityView, error)
	Create(ctx context.Context, c *Community) error
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, p UpdateParams) (*Community, error)
}

func NewRepository(db database.DB) Repository {
	return &postgresRepository{db: db}
}

type postgresRepository struct {
	db database.DB
}

func (r *postgresRepository) GetByName(ctx context.Context, name string) (*Community, error) {
	query := `
		SELECT id, name, owner_id, description, created_at, version
		FROM communities 
		WHERE name = $1`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, query, name)

	return scanCommunity(row)
}

func (r *postgresRepository) GetViewByName(
	ctx context.Context,
	name string,
) (*CommunityView, error) {
	query := `
		SELECT c.id, c.name, c.owner_id, c.description, c.created_at, c.version, u.username
		FROM communities c 
		LEFT JOIN users u ON u.id = c.owner_id
		WHERE c.name = $1`

	var c CommunityView
	var ownerID sql.NullInt64
	var username sql.NullString

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&c.ID, &c.Name, &ownerID,
		&c.Description, &c.CreatedAt, &c.Version, &username,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	if ownerID.Valid {
		c.OwnerID = &ownerID.Int64
		c.Owner = &CommunityOwner{
			Username: username.String,
		}
	}
	c.CreatedAt = c.CreatedAt.UTC()

	return &c, nil
}

func (r *postgresRepository) Create(ctx context.Context, c *Community) error {
	query := `
		INSERT INTO communities (name, owner_id, description)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, version`

	args := []any{c.Name, c.OwnerID, c.Description}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&c.ID, &c.CreatedAt, &c.Version)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "communities_name_key":
				return errs.ErrDuplicateCommunityName
			}
		}
		return err
	}

	c.CreatedAt = c.CreatedAt.UTC()
	return nil
}

func (r *postgresRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM communities WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *postgresRepository) Update(
	ctx context.Context,
	id int64,
	p UpdateParams,
) (*Community, error) {
	q := query.New()
	q.Update("communities")

	if p.Name != nil {
		q.Set("name = ?", *p.Name)
	}

	if p.Description != nil {
		q.Set("description = ?", *p.Description)
	}

	q.Set("version = version + 1")

	q.Where("id = ?", id)
	q.Returning("id, name, owner_id, description, created_at, version")
	query, args := q.ToSql()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, query, args...)
	community, err := scanCommunity(row)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "communities_name_key":
				return nil, errs.ErrDuplicateCommunityName
			}
		}
		return nil, err
	}

	return community, nil
}

func scanCommunity(row *sql.Row) (*Community, error) {
	var c Community
	var ownerID sql.NullInt64

	err := row.Scan(
		&c.ID, &c.Name, &ownerID,
		&c.Description, &c.CreatedAt, &c.Version,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	if ownerID.Valid {
		c.OwnerID = &ownerID.Int64
	}
	c.CreatedAt = c.CreatedAt.UTC()
	return &c, nil
}
