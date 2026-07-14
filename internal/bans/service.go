package bans

import (
	"context"

	"github.com/0xrinful/reddit-clone/internal/database"
)

type Service interface {
	Ban(ctx context.Context, actorID int64, p CreateParams) error
	Unban(ctx context.Context, actorID int64, p DeleteParams) error
	List(ctx context.Context, actorID int64, p ListParams) ([]*BanView, error)
}

type Authorizer interface {
	CanBan(ctx context.Context, communityID, actorID, targetID int64) error
	CanUnban(ctx context.Context, communityID, actorID, targetID int64) error
	CanViewBans(ctx context.Context, communityID, actorID int64) error
}

type UsersRepo interface {
	GetIDByUsername(ctx context.Context, username string) (int64, error)
}

type service struct {
	txBeginner database.TxBeginner
	authz      Authorizer
	bansRepo   Repository
	usersRepo  UsersRepo
}

func NewService(
	txBeginner database.TxBeginner,
	authz Authorizer,
	bansRepo Repository,
	usersRepo UsersRepo,
) Service {
	return &service{
		txBeginner: txBeginner,
		authz:      authz,
		bansRepo:   bansRepo,
		usersRepo:  usersRepo,
	}
}

func (s *service) Ban(ctx context.Context, actorID int64, p CreateParams) error {
	tx, err := s.txBeginner.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	ctxTx := database.WithTx(ctx, tx)

	targetID, err := s.usersRepo.GetIDByUsername(ctxTx, p.Username)
	if err != nil {
		return err
	}

	if err := s.authz.CanBan(ctxTx, p.CommunityID, actorID, targetID); err != nil {
		return err
	}

	b := &BanRecord{
		CommunityID: p.CommunityID,
		UserID:      targetID,
		BannedBy:    &actorID,
		Reason:      p.Reason,
		ExpiresAt:   p.Duration.Expiry(),
	}

	if err = s.bansRepo.Create(ctxTx, b); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *service) Unban(ctx context.Context, actorID int64, p DeleteParams) error {
	tx, err := s.txBeginner.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	ctxTx := database.WithTx(ctx, tx)

	targetID, err := s.usersRepo.GetIDByUsername(ctxTx, p.Username)
	if err != nil {
		return err
	}

	if err := s.authz.CanUnban(ctxTx, p.CommunityID, actorID, targetID); err != nil {
		return err
	}

	err = s.bansRepo.Delete(ctxTx, p.CommunityID, targetID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *service) List(ctx context.Context, actorID int64, p ListParams) ([]*BanView, error) {
	if err := s.authz.CanViewBans(ctx, p.CommunityID, actorID); err != nil {
		return nil, err
	}
	return s.bansRepo.List(ctx, p)
}
