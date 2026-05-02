package tokens

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/0xrinful/reddit-clone/internal/config"
	"github.com/0xrinful/reddit-clone/internal/shared/errs"
)

type Service interface {
	CreateRefreshToken(ctx context.Context, userID int64) (string, error)
	CreateAccessToken(userID int64) (string, error)
	VerifyAccessToken(tokenStr string) (*AccessTokenClaims, error)
	VerifyRefreshToken(ctx context.Context, plain string) (*Token, error)
	RotateRefreshToken(ctx context.Context, plain string, userID int64) (string, error)
	RevokeRefreshToken(ctx context.Context, plain string) error
}

func NewService(tokenRepo Repository, db *sql.DB, cfg config.JWTConfig) Service {
	return &service{
		db:         db,
		tokenRepo:  tokenRepo,
		secret:     []byte(cfg.Secret),
		refreshTTL: cfg.RefreshTokenTTL,
		accessTTL:  cfg.AccessTokenTTL,
	}
}

type service struct {
	db         *sql.DB
	tokenRepo  Repository
	secret     []byte
	refreshTTL time.Duration
	accessTTL  time.Duration
}

type AccessTokenClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *service) CreateAccessToken(userID int64) (string, error) {
	claims := AccessTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *service) VerifyAccessToken(tokenStr string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&AccessTokenClaims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errs.ErrInvalidToken
			}
			return s.secret, nil
		},
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
	)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired),
			errors.Is(err, jwt.ErrTokenNotValidYet),
			errors.Is(err, jwt.ErrTokenMalformed),
			errors.Is(err, jwt.ErrSignatureInvalid),
			errors.Is(err, jwt.ErrTokenUnverifiable):
			return nil, errs.ErrInvalidToken
		default:
			return nil, err
		}
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, errs.ErrInvalidToken
	}

	return claims, nil
}

func (s *service) CreateRefreshToken(ctx context.Context, userID int64) (string, error) {
	token, err := Generate(userID, s.refreshTTL, ScopeAuth)
	if err != nil {
		return "", err
	}

	err = s.tokenRepo.Insert(ctx, token)
	if err != nil {
		return "", err
	}

	return token.Plaintext, nil
}

func (s *service) VerifyRefreshToken(ctx context.Context, plain string) (*Token, error) {
	hash := Hash(plain)
	token, err := s.tokenRepo.GetByHash(ctx, ScopeAuth, hash)
	if err != nil {
		return nil, err
	}

	if time.Now().After(token.Expiry) {
		return nil, errs.ErrInvalidToken
	}

	return token, nil
}

func (s *service) RotateRefreshToken(
	ctx context.Context,
	plain string,
	userID int64,
) (string, error) {
	hash := Hash(plain)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	tokenRepo := NewRepository(tx)

	err = tokenRepo.DeleteByHash(ctx, hash)
	if err != nil {
		return "", err
	}

	newToken, err := Generate(userID, s.refreshTTL, ScopeAuth)
	if err != nil {
		return "", err
	}

	err = tokenRepo.Insert(ctx, newToken)
	if err != nil {
		return "", err
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}

	return newToken.Plaintext, nil
}

func (s *service) RevokeRefreshToken(ctx context.Context, plain string) error {
	hash := Hash(plain)
	return s.tokenRepo.DeleteByHash(ctx, hash)
}
