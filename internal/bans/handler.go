package bans

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

func NewHandler(srv Service, responder *response.Responder) *Handler {
	return &Handler{srv, responder}
}

func (h *Handler) Ban(w http.ResponseWriter, r *http.Request) {
	user, _ := request.GetUser(r)
	community := request.GetCommunity(r)

	var input CreateRequest

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

	params := CreateParams{
		UserID:      user.ID,
		CommunityID: community.ID,
		Username:    input.Username,
		Duration:    Duration(input.Duration),
		Reason:      input.Reason,
	}
	err = h.service.Create(r.Context(), params)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):

			// FIXME: change the response error for user not found
			h.responder.NotFound(w, r)
		case errors.Is(err, errs.ErrForbidden):
			h.responder.Forbidden(w)
		case errors.Is(err, errs.ErrSelfBan):

			// FIXME: change error response for self banning
			h.responder.Error(w, http.StatusBadRequest, "self ban", "can't ban self")
		default:
			h.responder.ServerError(w, err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UnBan(w http.ResponseWriter, r *http.Request) {
	user, _ := request.GetUser(r)
	community := request.GetCommunity(r)

	var input DeleteParams

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		h.responder.DecodeError(w, err)
		return
	}

	params := DeleteParams{
		UserID:      user.ID,
		CommunityID: community.ID,
		Username:    input.Username,
	}

	err = h.service.Delete(r.Context(), params)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrForbidden):
			h.responder.Forbidden(w)
		case errors.Is(err, errs.ErrNotFound):
			h.responder.NotFound(w, r)
		default:
			h.responder.ServerError(w, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
