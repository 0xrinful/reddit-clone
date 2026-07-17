package posts

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
	GetByID(ctx context.Context, id int64) (*Post, error)
	GetView(ctx context.Context, id int64) (*PostView, error)
	GetAuthorizationInfo(ctx context.Context, id int64) (*domain.PostAuthzInfo, error)
	Create(ctx context.Context, p *Post) error
	Update(ctx context.Context, id int64, p UpdateParams) (*Post, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, params ListParams) ([]*PostSummary, error)
	UpdateStatusByUser(ctx context.Context, userID, communityID int64,
		from, to domain.PostStatus) error
	ApplyScoreDelta(ctx context.Context, postID int64, delta int64) (int64, error)
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

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `
		SELECT 
			p.id, p.title, p.body, p.user_id, p.community_id, p.status, p.views, p.score, p.created_at 
		FROM posts p 
		WHERE p.id = $1`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := r.db(ctx).QueryRowContext(ctx, query, id)
	return scanPost(row)
}

func (r *postgresRepository) GetAuthorizationInfo(
	ctx context.Context,
	id int64,
) (*domain.PostAuthzInfo, error) {
	query := `
		SELECT user_id, community_id, status FROM posts 
		WHERE id = $1`

	var info domain.PostAuthzInfo

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.db(ctx).
		QueryRowContext(ctx, query, id).
		Scan(&info.AuthorID, &info.CommunityID, &info.Status)
	if err != nil {
		return nil, scanError(err)
	}

	return &info, nil
}

func (r *postgresRepository) GetView(
	ctx context.Context,
	id int64,
) (*PostView, error) {
	query := `
		SELECT 
			p.id, p.title, p.body, p.created_at, p.status, p.score, p.views,
			p.user_id, u.username as author,
			p.community_id, c.name 
		FROM posts p 
		JOIN communities c ON p.community_id = c.id
		JOIN users u ON p.user_id = u.id
		WHERE p.id = $1`

	var p PostView

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.db(ctx).QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Title, &p.Body, &p.CreatedAt, &p.Status, &p.Score, &p.Views,
		&p.UserID, &p.Author.Username,
		&p.CommunityID, &p.Community.Name,
	)
	if err != nil {
		return nil, scanError(err)
	}

	p.CreatedAt = p.CreatedAt.UTC()

	return &p, nil
}

func (r *postgresRepository) Create(ctx context.Context, p *Post) error {
	query := `
		INSERT INTO posts (title, body, user_id, community_id, score)
		VALUES ($1, $2, $3, $4, 1)
		RETURNING id, created_at`

	args := []any{p.Title, p.Body, p.UserID, p.CommunityID}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.db(ctx).QueryRowContext(ctx, query, args...).Scan(&p.ID, &p.CreatedAt)
	if err != nil {
		return err
	}
	p.CreatedAt = p.CreatedAt.UTC()
	p.Score = 1

	return nil
}

func (r *postgresRepository) Update(ctx context.Context, id int64, p UpdateParams) (*Post, error) {
	q := query.New()
	q.Update("posts")

	if p.Title != nil {
		q.Set("title = ?", *p.Title)
	}

	if p.Body != nil {
		q.Set("body = ?", *p.Body)
	}

	q.Where("id = ?", id)
	q.Returning("*")
	query, args := q.ToSql()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := r.db(ctx).QueryRowContext(ctx, query, args...)
	return scanPost(row)
}

func (r *postgresRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := r.db(ctx).ExecContext(ctx, query, id)
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

func (r *postgresRepository) UpdateStatusByUser(
	ctx context.Context,
	userID, communityID int64,
	from, to domain.PostStatus,
) error {
	query := `
		UPDATE posts SET status = $4
		WHERE user_id = $1 AND community_id = $2 AND status = $3`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := r.db(ctx).ExecContext(ctx, query, userID, communityID, from, to)
	return err
}

func (r *postgresRepository) List(
	ctx context.Context,
	params ListParams,
) ([]*PostSummary, error) {
	q := query.New()
	cursor := params.Pagination.Cursor

	q.Select("p.id, p.title, p.body, p.score, p.created_at, u.username, c.name", "posts p")
	q.Join("communities c ON p.community_id = c.id")
	q.Join("users u ON p.user_id = u.id")

	q.Where("p.community_id = ?", params.CommunityID)
	q.Where("p.status = ?", domain.PostStatusActive)

	switch params.Sort {
	case SortByNew:
		if cursor != nil {
			q.Where("(p.created_at, p.id) < (?, ?)", cursor.CreatedAt, cursor.ID)
		}
		q.Order("p.created_at DESC, p.id DESC")
	case SortByTop, SortByHot:
		if cursor != nil {
			q.Where("(p.score, p.id) < (?, ?)", cursor.Score, cursor.ID)
		}
		q.Order("p.score DESC, p.id DESC")
	default:
		panic("invalid sort value")
	}

	q.Limit(params.Pagination.Limit)
	query, args := q.ToSql()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := r.db(ctx).QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*PostSummary, 0, params.Pagination.Limit)
	for rows.Next() {
		var p PostSummary
		err = rows.Scan(
			&p.ID, &p.Title, &p.Body, &p.Score, &p.CreatedAt,
			&p.Author.Username, &p.Community.Name,
		)
		if err != nil {
			return nil, err
		}
		p.CreatedAt = p.CreatedAt.UTC()
		posts = append(posts, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *postgresRepository) ApplyScoreDelta(
	ctx context.Context,
	postID int64,
	delta int64,
) (int64, error) {
	query := `
		UPDATE posts SET score = score + $1
		WHERE id = $2
		RETURNING score`

	var newScore int64

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.db(ctx).QueryRowContext(ctx, query, delta, postID).Scan(&newScore)
	if err != nil {
		return 0, scanError(err)
	}

	return newScore, nil
}

func scanPost(row *sql.Row) (*Post, error) {
	var p Post

	err := row.Scan(
		&p.ID, &p.Title, &p.Body, &p.UserID, &p.CommunityID,
		&p.Status, &p.Views, &p.Score, &p.CreatedAt,
	)
	if err != nil {
		return nil, scanError(err)
	}

	p.CreatedAt = p.CreatedAt.UTC()

	return &p, nil
}

func scanError(err error) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return errs.ErrNotFound
	}
	return err
}
