package bans

import (
	"context"

	"github.com/0xrinful/reddit-clone/internal/database"
	"github.com/0xrinful/reddit-clone/internal/domain"
	"github.com/0xrinful/reddit-clone/internal/shared/errs"
)

type Service interface {
	Ban(ctx context.Context, requesterID int64, p CreateParams) error
	Unban(ctx context.Context, requesterID int64, p DeleteParams) error
}

type MembersRepo interface {
	GetRole(ctx context.Context, communityID, userID int64) (domain.Role, error)
}

type UsersRepo interface {
	GetIDByUsername(ctx context.Context, username string) (int64, error)
}

type service struct {
	txBeginner  database.TxBeginner
	bansRepo    Repository
	membersRepo MembersRepo
	usersRepo   UsersRepo
}

func NewService(
	txBeginner database.TxBeginner,
	bansRepo Repository,
	membersRepo MembersRepo,
	usersRepo UsersRepo,
) Service {
	return &service{
		txBeginner:  txBeginner,
		bansRepo:    bansRepo,
		membersRepo: membersRepo,
		usersRepo:   usersRepo,
	}
}

func (s *service) Ban(ctx context.Context, requesterID int64, p CreateParams) error {
	tx, err := s.txBeginner.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	ctxTx := database.WithTx(ctx, tx)

	requesterRole, err := s.membersRepo.GetRole(ctxTx, p.CommunityID, requesterID)
	if err != nil {
		return err
	}

	if !requesterRole.IsOwner() && !requesterRole.IsModerator() {
		return errs.ErrPermissionDenied
	}

	targetID, err := s.usersRepo.GetIDByUsername(ctxTx, p.Username)
	if err != nil {
		return err
	}

	if requesterID == targetID {
		return errs.ErrSelfBan
	}

	targetRole, err := s.membersRepo.GetRole(ctxTx, p.CommunityID, targetID)
	if err != nil {
		return err
	}

	if !requesterRole.CanManageBan(targetRole) {
		return errs.ErrPermissionDenied
	}

	b := &BanRecord{
		CommunityID: p.CommunityID,
		UserID:      targetID,
		BannedBy:    &requesterID,
		Reason:      p.Reason,
		ExpiresAt:   p.Duration.Expiry(),
	}

	if err = s.bansRepo.Create(ctxTx, b); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *service) Unban(ctx context.Context, requesterID int64, p DeleteParams) error {
	tx, err := s.txBeginner.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	ctxTx := database.WithTx(ctx, tx)

	requesterRole, err := s.membersRepo.GetRole(ctxTx, p.CommunityID, requesterID)
	if err != nil {
		return err
	}

	if !requesterRole.IsOwner() && !requesterRole.IsModerator() {
		return errs.ErrPermissionDenied
	}

	targetID, err := s.usersRepo.GetIDByUsername(ctxTx, p.Username)
	if err != nil {
		return err
	}

	targetRole, err := s.membersRepo.GetRole(ctxTx, p.CommunityID, targetID)
	if err != nil {
		return err
	}

	if !requesterRole.CanManageBan(targetRole) {
		return errs.ErrPermissionDenied
	}

	err = s.bansRepo.Delete(ctxTx, p.CommunityID, targetID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
