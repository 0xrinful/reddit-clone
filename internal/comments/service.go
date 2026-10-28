package comments

import (
	"context"

	"github.com/0xrinful/reddit-clone/internal/database"
	"github.com/0xrinful/reddit-clone/internal/domain"
)

type Service interface {
	Create(ctx context.Context, params CreateParams) (*Comment, error)
	Reply(ctx context.Context, params ReplyParams) (*Comment, error)
	ListByPost(ctx context.Context, params ListByPostParams) ([]*CommentView, error)
}

type Authorizer interface {
	CanComment(ctx context.Context, actorID, postID int64) error
	CanReply(ctx context.Context, actorID, postID int64, status domain.CommentStatus) error
}

func NewService(
	txBeginner database.TxBeginner,
	authz Authorizer,
	commentsRepo Repository,
) Service {
	return &service{
		txBeginner:   txBeginner,
		authz:        authz,
		commentsRepo: commentsRepo,
	}
}

type service struct {
	txBeginner   database.TxBeginner
	authz        Authorizer
	commentsRepo Repository
}

func (s *service) Create(ctx context.Context, params CreateParams) (*Comment, error) {
	tx, err := s.txBeginner.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	ctxTx := database.WithTx(ctx, tx)

	if err := s.authz.CanComment(ctx, params.UserID, params.PostID); err != nil {
		return nil, err
	}

	c := &Comment{
		PostID:  params.PostID,
		UserID:  params.UserID,
		Content: params.Content,
	}

	if err := s.commentsRepo.Create(ctxTx, c); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return c, nil
}

func (s *service) Reply(ctx context.Context, params ReplyParams) (*Comment, error) {
	tx, err := s.txBeginner.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	ctxTx := database.WithTx(ctx, tx)

	parent, err := s.commentsRepo.GetByID(ctx, params.ParentID)
	if err != nil {
		return nil, err
	}

	if err := s.authz.CanReply(ctx, params.UserID, parent.PostID, parent.Status); err != nil {
		return nil, err
	}

	c := &Comment{
		PostID:   parent.PostID,
		UserID:   params.UserID,
		ParentID: &params.ParentID,
		Content:  params.Content,
	}

	if err := s.commentsRepo.Create(ctxTx, c); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return c, nil
}

func (s *service) ListByPost(ctx context.Context, params ListByPostParams) ([]*CommentView, error) {
	return s.commentsRepo.ListByPost(ctx, params)
}
