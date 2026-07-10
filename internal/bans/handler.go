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
		CommunityID: community.ID,
		Username:    input.Username,
		Reason:      input.Reason,
		Duration:    input.Duration,
	}

	err = h.service.Ban(r.Context(), user.ID, params)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			h.responder.NotFound(w, r)
		case errors.Is(err, errs.ErrPermissionDenied):
			h.responder.PermissionDenied(w)
		case errors.Is(err, errs.ErrSelfBan):
			h.responder.ForbiddenMsg(w, "self_ban", "can't ban self")
		default:
			h.responder.ServerError(w, err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Unban(w http.ResponseWriter, r *http.Request) {
	user, _ := request.GetUser(r)
	community := request.GetCommunity(r)

	params := DeleteParams{
		CommunityID: community.ID,
		Username:    r.PathValue("username"),
	}

	err := h.service.Unban(r.Context(), user.ID, params)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			h.responder.NotFound(w, r)
		case errors.Is(err, errs.ErrPermissionDenied):
			h.responder.Forbidden(w)
		default:
			h.responder.ServerError(w, err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
