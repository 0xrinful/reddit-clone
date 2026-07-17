package posts

import (
	"context"

	"github.com/0xrinful/reddit-clone/internal/database"
	"github.com/0xrinful/reddit-clone/internal/domain"
)

type Service interface {
	Get(ctx context.Context, id int64, actorID *int64) (*PostView, error)
	Create(ctx context.Context, params CreateParams) (*Post, error)
	Update(ctx context.Context, id, actorID int64, p UpdateParams) (*Post, error)
	Delete(ctx context.Context, id, userID int64) error
	List(ctx context.Context, params ListParams) ([]*PostSummary, error)
}

type Authorizer interface {
	CanPost(ctx context.Context, communityID, userID int64) error
	CanViewPost(ctx context.Context, communityID int64, actorID *int64, authorID int64,
		status domain.PostStatus) error
	CanUpdatePost(ctx context.Context, communityID, actorID, authorID int64) error
	CanDeletePost(actorID, authorID int64) error
}

type VotesRepo interface {
	CreateInitialPostVote(ctx context.Context, userID, postID int64) error
}

func NewService(
	txBeginner database.TxBeginner,
	authz Authorizer,
	postsRepo Repository,
	votesRepo VotesRepo,
) Service {
	return &service{
		txBeginner: txBeginner,
		authz:      authz,
		postsRepo:  postsRepo,
		votesRepo:  votesRepo,
	}
}

type service struct {
	txBeginner database.TxBeginner
	authz      Authorizer
	postsRepo  Repository
	votesRepo  VotesRepo
}

func (s *service) Create(ctx context.Context, params CreateParams) (*Post, error) {
	tx, err := s.txBeginner.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	ctxTx := database.WithTx(ctx, tx)

	if err := s.authz.CanPost(ctxTx, params.CommunityID, params.UserID); err != nil {
		return nil, err
	}

	p := &Post{
		Title:       params.Title,
		Body:        params.Body,
		UserID:      params.UserID,
		CommunityID: params.CommunityID,
	}

	if err := s.postsRepo.Create(ctxTx, p); err != nil {
		return nil, err
	}

	if err := s.votesRepo.CreateInitialPostVote(ctxTx, params.UserID, p.ID); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *service) Update(ctx context.Context, id, actorID int64, p UpdateParams) (*Post, error) {
	post, err := s.postsRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.authz.CanUpdatePost(ctx, post.CommunityID, actorID, post.UserID); err != nil {
		return nil, err
	}

	return s.postsRepo.Update(ctx, id, p)
}

func (s *service) Get(ctx context.Context, id int64, actorID *int64) (*PostView, error) {
	post, err := s.postsRepo.GetView(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = s.authz.CanViewPost(ctx, post.CommunityID, actorID, post.UserID,
		post.Status); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *service) Delete(ctx context.Context, id, actorID int64) error {
	post, err := s.postsRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.authz.CanDeletePost(actorID, post.UserID); err != nil {
		return err
	}

	return s.postsRepo.Delete(ctx, id)
}

func (s *service) List(
	ctx context.Context,
	params ListParams,
) ([]*PostSummary, error) {
	return s.postsRepo.List(ctx, params)
}
