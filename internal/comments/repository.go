package comments

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
	GetByID(ctx context.Context, id int64) (*Comment, error)
	Create(ctx context.Context, c *Comment) error
	ListByPost(ctx context.Context, params ListByPostParams) ([]*CommentView, error)
	ListReplies(ctx context.Context, params ListRepliesParams) ([]*CommentView, error)
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
		RETURNING id, score, reply_count, created_at`

	args := []any{c.PostID, c.UserID, c.ParentID, c.Content}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.db(ctx).
		QueryRowContext(ctx, query, args...).
		Scan(&c.ID, &c.Score, &c.ReplyCount, &c.CreatedAt)
	if err != nil {
		return err
	}
	c.CreatedAt = c.CreatedAt.UTC()

	return nil
}

func (r *postgresRepository) ListByPost(
	ctx context.Context,
	params ListByPostParams,
) ([]*CommentView, error) {
	q := query.New()
	cursor := params.Pagination.Cursor

	q.Select(
		"c.id, c.post_id, c.user_id, u.username, c.parent_id, c.content, c.score, c.reply_count, c.status, c.created_at",
		"comments c",
	)
	q.Join("users u ON c.user_id = u.id")

	q.Where("c.depth = 0")
	q.Where("c.post_id = ?", params.PostID)
	q.Where("c.status = ?", domain.CommentStatusVisible)

	switch params.Sort {
	case SortByNew:
		if cursor != nil {
			q.Where("(c.created_at , c.id) < (?, ?)", cursor.CreatedAt, cursor.ID)
		}
		q.Order("c.created_at DESC, c.id DESC")
	case SortByTop:
		if cursor != nil {
			q.Where("(c.score , c.id) < (?, ?)", cursor.Score, cursor.ID)
		}
		q.Order("c.score DESC, c.id DESC")
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

	comments := make([]*CommentView, 0, params.Pagination.Limit)
	for rows.Next() {
		c, err := scanView(rows)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return comments, nil
}

func (r *postgresRepository) ListReplies(
	ctx context.Context,
	params ListRepliesParams,
) ([]*CommentView, error) {
	q := query.New()
	cursor := params.Pagination.Cursor

	q.Select(
		"c.id, c.post_id, c.user_id, u.username, c.parent_id, c.content, c.score, c.reply_count, c.status, c.created_at",
		"comments c",
	)
	q.Join("users u ON c.user_id = u.id")

	q.Where("c.parent_id = ?", params.CommentID)
	q.Where("c.status = ?", domain.CommentStatusVisible)

	if cursor != nil {
		q.Where("(c.created_at, c.id) > (?, ?)", cursor.CreatedAt, cursor.ID)
	}
	q.Order("c.created_at ASC, c.id ASC")

	q.Limit(params.Pagination.Limit)
	query, args := q.ToSql()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := r.db(ctx).QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]*CommentView, 0, params.Pagination.Limit)
	for rows.Next() {
		c, err := scanView(rows)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return comments, nil
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

func scanView(rows *sql.Rows) (*CommentView, error) {
	var c CommentView
	var parentID sql.NullInt64

	err := rows.Scan(&c.ID, &c.PostID,
		&c.UserID, &c.Author.Username, &parentID,
		&c.Content, &c.Score, &c.ReplyCount,
		&c.Status, &c.CreatedAt)
	if err != nil {
		return nil, err
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
