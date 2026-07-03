package auth

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/0xrinful/reddit-clone/internal/database"
	"github.com/0xrinful/reddit-clone/internal/shared/background"
	"github.com/0xrinful/reddit-clone/internal/shared/errs"
	"github.com/0xrinful/reddit-clone/internal/shared/mailer"
	"github.com/0xrinful/reddit-clone/internal/tokens"
	"github.com/0xrinful/reddit-clone/internal/users"
)

type Service interface {
	RegisterUser(ctx context.Context, params CreateUserParams) (*users.User, error)
	SendActivationEmail(ctx context.Context, email string) error
	ActivateUser(ctx context.Context, plaintext string) error
	AuthenticateUser(ctx context.Context, email string, password string) (*users.User, error)
}

func NewService(
	txBeginner database.TxBeginner,
	userRepo users.Repository,
	tokenRepo tokens.Repository,
	mailer *mailer.Mailer,
	logger *slog.Logger,
	background *background.Worker,
) *service {
	return &service{
		txBeginner: txBeginner,
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		mailer:     mailer,
		logger:     logger,
		background: background,
	}
}

type service struct {
	txBeginner database.TxBeginner
	userRepo   users.Repository
	tokenRepo  tokens.Repository
	mailer     *mailer.Mailer
	logger     *slog.Logger
	background *background.Worker
}

type CreateUserParams struct {
	Username      string
	Email         string
	PlainPassword string
}

func (s *service) RegisterUser(ctx context.Context, params CreateUserParams) (*users.User, error) {
	user := &users.User{
		Username: params.Username,
		Email:    params.Email,
	}

	err := user.Password.Set(params.PlainPassword)
	if err != nil {
		return nil, err
	}

	tx, err := s.txBeginner.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	ctxTx := database.WithTx(ctx, tx)

	err = s.userRepo.Create(ctxTx, user)
	if err != nil {
		return nil, err
	}

	token, err := tokens.Generate(user.ID, 24*time.Hour, tokens.ScopeActivation)
	if err != nil {
		return nil, err
	}

	if err = s.tokenRepo.Insert(ctxTx, token); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	s.background.Run(func() {
		s.deliverActivationEmail(user.ID, user.Email, user.Username, token.Plaintext)
	})

	return user, nil
}

func (s *service) ActivateUser(ctx context.Context, plaintext string) error {
	hashed := tokens.Hash(plaintext)
	token, err := s.tokenRepo.GetByHash(ctx, tokens.ScopeActivation, hashed)
	if err != nil {
		return err
	}

	if time.Now().After(token.Expiry) {
		return errs.ErrInvalidToken
	}

	err = s.userRepo.SetActivated(ctx, token.UserID)
	if err != nil {
		return err
	}

	if err = s.tokenRepo.DeleteAllForUser(ctx, tokens.ScopeActivation, token.UserID); err != nil {
		s.logger.Error("failed to delete activation tokens", "userID", token.UserID, "error", err)
	}

	return nil
}

func (s *service) SendActivationEmail(ctx context.Context, email string) error {
	tx, err := s.txBeginner.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ctxTx := database.WithTx(ctx, tx)

	user, err := s.userRepo.GetByEmail(ctxTx, email)
	if err != nil {
		return err
	}

	if user.Activated {
		return errs.ErrAlreadyActivated
	}

	token, err := tokens.Generate(user.ID, 24*time.Hour, tokens.ScopeActivation)
	if err != nil {
		return err
	}

	if err = s.tokenRepo.DeleteAllForUser(ctxTx, tokens.ScopeActivation, user.ID); err != nil {
		return err
	}

	if err = s.tokenRepo.Insert(ctxTx, token); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	s.background.Run(func() {
		s.deliverActivationEmail(user.ID, user.Email, user.Username, token.Plaintext)
	})

	return nil
}

func (s *service) deliverActivationEmail(userID int64, email, username, plaintext string) {
	err := s.mailer.Send(email, "activation.html", map[string]any{
		"Username":      username,
		"ActivationURL": "https://app.com/verify?token=" + plaintext, // TODO: update url
	})
	if err != nil {
		s.logger.Error("activation email failed", "userID", userID, "error", err)
	}
}

func (s *service) AuthenticateUser(
	ctx context.Context,
	email string,
	password string,
) (*users.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if errors.Is(err, errs.ErrNotFound) {
		return nil, errs.ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	match, err := user.Password.Match(password)
	if err != nil {
		return nil, err
	}

	if !match {
		return nil, errs.ErrInvalidCredentials
	}

	return user, nil
}
