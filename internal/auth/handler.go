package auth

import (
	"errors"
	"net/http"

	"github.com/0xrinful/reddit-clone/internal/shared/errs"
	"github.com/0xrinful/reddit-clone/internal/shared/request"
	"github.com/0xrinful/reddit-clone/internal/shared/response"
	"github.com/0xrinful/reddit-clone/internal/shared/validator"
)

type Handler struct {
	service   Service
	responder *response.Responder
}

func NewHandler(svc Service, responder *response.Responder) *Handler {
	return &Handler{svc, responder}
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

	user, err := h.service.RegisterUser(r.Context(), params)
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

	err = h.service.ActivateUser(r.Context(), input.Token)
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
