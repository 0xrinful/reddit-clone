package auth

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

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
}

func NewService(
	db *sql.DB,
	userRepo users.Repository,
	tokenRepo tokens.Repository,
	mailer *mailer.Mailer,
	logger *slog.Logger,
	background *background.Worker,
) *service {
	return &service{
		db:         db,
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		mailer:     mailer,
		logger:     logger,
		background: background,
	}
}

type service struct {
	db         *sql.DB
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

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	userRepo := users.NewRepository(tx)
	tokenRepo := tokens.NewRepository(tx)

	err = userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	token, err := tokens.Generate(user.ID, 24*time.Hour, tokens.ScopeActivation)
	if err != nil {
		return nil, err
	}

	if err = tokenRepo.Insert(ctx, token); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	s.background.Run(func() {
		s.sendActivation(user.ID, user.Email, user.Username, token.Plaintext)
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
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	userRepo := users.NewRepository(tx)
	tokenRepo := tokens.NewRepository(tx)

	user, err := userRepo.GetByEmail(ctx, email)
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

	if err = tokenRepo.DeleteAllForUser(ctx, tokens.ScopeActivation, user.ID); err != nil {
		return err
	}

	if err = tokenRepo.Insert(ctx, token); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	s.background.Run(func() {
		s.sendActivation(user.ID, user.Email, user.Username, token.Plaintext)
	})

	return nil
}

func (s *service) sendActivation(userID int64, email, username, plaintext string) {
	err := s.mailer.Send(email, "activation.html", map[string]any{
		"Username":      username,
		"ActivationURL": "https://app.com/verify?token=" + plaintext,
	})
	if err != nil {
		s.logger.Error("activation email failed", "userID", userID, "error", err)
	}
}
