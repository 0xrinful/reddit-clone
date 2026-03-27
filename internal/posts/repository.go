package posts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/0xrinful/reddit-clone/internal/shared/apperr"
	"github.com/0xrinful/reddit-clone/internal/shared/query"
)

type Repository interface {
	Get(ctx context.Context, id, communityID int64) (*Post, error)
	Create(ctx context.Context, p *Post) error
	Update(ctx context.Context, p UpdatePostParams) error
	Delete(ctx context.Context, id, userID, communityID int64) error
	List(ctx context.Context, params ListPostParams) ([]*Post, error)
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{db}
}

type postgresRepository struct {
	db *sql.DB
}

func (r *postgresRepository) Get(ctx context.Context, id, CommunityID int64) (*Post, error) {
	query := `
		SELECT id, title, body, user_id, community_id, views, created_at, version
		FROM posts
		WHERE id = $1 AND community_id = $2`

	var p Post

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, id, CommunityID).Scan(
		&p.ID, &p.Title, &p.Body, &p.UserID, &p.CommunityID,
		&p.Views, &p.CreatedAt, &p.Version,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperr.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *postgresRepository) Create(ctx context.Context, p *Post) error {
	query := `
		INSERT INTO posts (title, body, user_id, community_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []any{p.Title, p.Body, p.UserID, p.CommunityID}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return r.db.QueryRowContext(ctx, query, args...).Scan(&p.ID, &p.CreatedAt, &p.Version)
}

func (r *postgresRepository) Update(ctx context.Context, p UpdatePostParams) error {
	var q query.Query
	q.Update("posts")

	if p.Title != nil {
		q.Set("title", *p.Title)
	}

	if p.Body != nil {
		q.Set("body", *p.Body)
	}

	q.Set("version", query.Raw("version + 1"))

	q.Where("id = ? AND community_id = ? AND user_id = ?", p.ID, p.CommunityID, p.UserID)
	query, args := q.ToSql()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return apperr.ErrNotFound
	}

	return nil
}

func (r *postgresRepository) Delete(ctx context.Context, id, userID, communityID int64) error {
	query := `
		DELETE FROM posts 
		WHERE id = $1 AND user_id = $2 AND community_id = $3`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, id, userID, communityID)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return apperr.ErrNotFound
	}

	return nil
}

func (r *postgresRepository) List(ctx context.Context, params ListPostParams) ([]*Post, error) {
	query := fmt.Sprintf(`
		SELECT p.id, p.title, p.body, p.user_id, p.community_id, p.views, p.created_at,
		p.version, COALESCE(SUM(v.value), 0) AS score FROM posts p 
		LEFT JOIN post_votes v ON p.id = v.post_id
		WHERE community_id = $1 AND ($2 = 0 OR p.id < $2)
		GROUP BY p.id
		ORDER BY %s, p.id DESC
		LIMIT $3`, params.Sort.ToSql())

	args := []any{params.CommunityID, params.Cursor.After, params.Cursor.Limit}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*Post, 0, params.Cursor.Limit)
	for rows.Next() {
		var p Post
		err = rows.Scan(
			&p.ID, &p.Title, &p.Body, &p.UserID, &p.CommunityID,
			&p.Views, &p.CreatedAt, &p.Version, &p.Score,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
