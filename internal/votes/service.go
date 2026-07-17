package votes

import (
	"context"

	"github.com/0xrinful/reddit-clone/internal/database"
)

type Service interface {
	VotePost(ctx context.Context, vote PostVote) (*VotePostResult, error)
}

type Authorizer interface {
	CanVotePost(ctx context.Context, actorID, postID int64) error
}

type PostsRepo interface {
	ApplyScoreDelta(ctx context.Context, postID int64, delta int64) (int64, error)
}

func NewService(
	txBeginner database.TxBeginner,
	authz Authorizer,
	votesRepo Repository,
	postsRepo PostsRepo,
) Service {
	return &service{
		txBeginner: txBeginner,
		authz:      authz,
		votesRepo:  votesRepo,
		postsRepo:  postsRepo,
	}
}

type service struct {
	txBeginner database.TxBeginner
	authz      Authorizer
	votesRepo  Repository
	postsRepo  PostsRepo
}

func (s *service) VotePost(ctx context.Context, vote PostVote) (*VotePostResult, error) {
	tx, err := s.txBeginner.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	ctxTx := database.WithTx(ctx, tx)

	if err := s.authz.CanVotePost(ctxTx, vote.UserID, vote.PostID); err != nil {
		return nil, err
	}

	delta, err := s.votesRepo.VotePost(ctxTx, vote)
	if err != nil {
		return nil, err
	}

	newScore, err := s.postsRepo.ApplyScoreDelta(ctxTx, vote.PostID, delta)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	result := &VotePostResult{
		PostID: vote.PostID,
		Score:  newScore,
		Value:  vote.Value,
	}

	return result, nil
}
