package bans

import (
	"context"

	"github.com/0xrinful/reddit-clone/internal/database"
	"github.com/0xrinful/reddit-clone/internal/members"
	"github.com/0xrinful/reddit-clone/internal/shared/errs"
	"github.com/0xrinful/reddit-clone/internal/users"
)

type Service interface {
	Create(ctx context.Context, p CreateParams) error
	Delete(ctx context.Context, p DeleteParams) error
}

type service struct {
	txBeginner  database.TxBeginner
	modsRepo    Repository
	membersRepo members.Repository
	usersRepo   users.Repository
}

func NewService(
	txBeginner database.TxBeginner,
	bansRepo Repository,
	membersRepo members.Repository,
	usersRepo users.Repository,
) Service {
	return &service{
		txBeginner:  txBeginner,
		modsRepo:    bansRepo,
		membersRepo: membersRepo,
		usersRepo:   usersRepo,
	}
}

func (s *service) Create(ctx context.Context, p CreateParams) error {
	tx, err := s.txBeginner.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	ctxTx := database.WithTx(ctx, tx)

	u, err := s.membersRepo.Get(ctxTx, p.CommunityID, p.UserID)
	if err != nil {
		return err
	}

	if u.Role != members.RoleModerator && u.Role != members.RoleOwner {
		return errs.ErrForbidden
	}

	t, err := s.usersRepo.GetByUsername(ctxTx, p.Username)
	if err != nil {
		return err
	}

	// NOTE: overrules moderator self banning and owner banning self
	if t.ID == u.UserID {
		return errs.ErrSelfBan
	}

	// TODO: maybe check if moderator banning owner?

	b := &BanRecord{
		BannedBy:    p.UserID,
		CommunityID: p.CommunityID,
		UserID:      t.ID,
		Reason:      p.Reason,
		Expiry:      p.Duration.Expiry(),
	}

	if err = s.modsRepo.Create(ctxTx, b); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *service) Delete(ctx context.Context, p DeleteParams) error {
	tx, err := s.txBeginner.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	ctxTx := database.WithTx(ctx, tx)

	u, err := s.membersRepo.Get(ctxTx, p.CommunityID, p.UserID)
	if err != nil {
		return err
	}

	if u.Role != members.RoleModerator && u.Role != members.RoleOwner {
		return errs.ErrForbidden
	}

	t, err := s.usersRepo.GetByUsername(ctxTx, p.Username)
	if err != nil {
		return err
	}

	err = s.modsRepo.Delete(ctx, p.CommunityID, t.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
