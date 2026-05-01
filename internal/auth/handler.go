package auth

import (
	"errors"
	"net/http"

	"github.com/0xrinful/reddit-clone/internal/shared/errs"
	"github.com/0xrinful/reddit-clone/internal/shared/request"
	"github.com/0xrinful/reddit-clone/internal/shared/response"
	"github.com/0xrinful/reddit-clone/internal/shared/validator"
	"github.com/0xrinful/reddit-clone/internal/tokens"
)

type Handler struct {
	authService   Service
	tokensService tokens.Service
	responder     *response.Responder
}

func NewHandler(
	authSvc Service,
	tokensSvc tokens.Service,
	responder *response.Responder,
) *Handler {
	return &Handler{authSvc, tokensSvc, responder}
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var input RegisterUserRequest

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		h.responder.DecodeError(w, err)
		return
	}

	v := validator.New()
	if input.Validate(v); !v.Valid() {
		h.responder.ValidationError(w, v.Errors)
		return
	}

	params := CreateUserParams{
		Username:      input.Username,
		Email:         input.Email,
		PlainPassword: input.Password,
	}

	user, err := h.authService.RegisterUser(r.Context(), params)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrDuplicateEmail):
			v.AddError("email", "email is already in use")
			h.responder.ValidationError(w, v.Errors)
		case errors.Is(err, errs.ErrDuplicateUsername):
			v.AddError("username", "username is already in use")
			h.responder.ValidationError(w, v.Errors)
		default:
			h.responder.ServerError(w, err)
		}
		return
	}

	h.responder.JSON(w, http.StatusCreated, toRegisterUserResponse(user))
}

func (h *Handler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	var input ActivateUserRequest

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		h.responder.DecodeError(w, err)
		return
	}

	v := validator.New()
	if input.Validate(v); !v.Valid() {
		h.responder.ValidationError(w, v.Errors)
		return
	}

	err = h.authService.ActivateUser(r.Context(), input.Token)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidToken):
			v.AddError("token", "invalid or expired token")
			h.responder.ValidationError(w, v.Errors)
		default:
			h.responder.ServerError(w, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) SendActivationEmail(w http.ResponseWriter, r *http.Request) {
	var input SendActivationEmailRequest

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		h.responder.DecodeError(w, err)
		return
	}

	v := validator.New()
	if input.Validate(v); !v.Valid() {
		h.responder.ValidationError(w, v.Errors)
		return
	}

	err = h.authService.SendActivationEmail(r.Context(), input.Email)
	if err != nil && !errors.Is(err, errs.ErrNotFound) &&
		!errors.Is(err, errs.ErrAlreadyActivated) {
		h.responder.ServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input LoginRequest

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		h.responder.DecodeError(w, err)
		return
	}

	v := validator.New()
	if input.Validate(v); !v.Valid() {
		h.responder.ValidationError(w, v.Errors)
		return
	}

	user, err := h.authService.AuthenticateUser(r.Context(), input.Email, input.Password)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidCredentials):
			h.responder.InvalidCredentials(w)
		default:
			h.responder.ServerError(w, err)
		}
		return
	}

	refreshToken, err := h.tokensService.CreateRefreshToken(r.Context(), user.ID)
	if err != nil {
		h.responder.ServerError(w, err)
		return
	}

	accessToken, err := h.tokensService.CreateAccessToken(user.ID)
	if err != nil {
		h.responder.ServerError(w, err)
		return
	}

	h.responder.JSON(w, http.StatusOK, toLoginResponse(accessToken, refreshToken))
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var input RefreshRequest

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		h.responder.DecodeError(w, err)
		return
	}

	v := validator.New()
	if input.Validate(v); !v.Valid() {
		h.responder.ValidationError(w, v.Errors)
		return
	}

	oldToken, err := h.tokensService.VerifyRefreshToken(r.Context(), input.Token)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidToken):
			h.responder.InvalidToken(w)
		default:
			h.responder.ServerError(w, err)
		}
		return
	}

	refreshToken, err := h.tokensService.RotateRefreshToken(
		r.Context(),
		input.Token,
		oldToken.UserID,
	)
	if err != nil {
		h.responder.ServerError(w, err)
		return
	}

	accessToken, err := h.tokensService.CreateAccessToken(oldToken.UserID)
	if err != nil {
		h.responder.ServerError(w, err)
		return
	}

	h.responder.JSON(w, http.StatusOK, toRefreshResponse(accessToken, refreshToken))
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var input LogoutRequest

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		h.responder.DecodeError(w, err)
		return
	}

	v := validator.New()
	if input.Validate(v); !v.Valid() {
		h.responder.ValidationError(w, v.Errors)
		return
	}

	err = h.tokensService.RevokeRefreshToken(r.Context(), input.Token)
	if err != nil {
		h.responder.ServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
