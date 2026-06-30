package communities

import (
	"errors"
	"fmt"
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

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("community_name")

	community, err := h.service.GetViewByName(r.Context(), name)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			h.responder.NotFound(w, r)
		default:
			h.responder.ServerError(w, err)
		}
		return
	}

	h.responder.JSON(w, http.StatusOK, toCommunityResponse(community))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user, _ := request.GetUser(r)

	var input CreateCommunitytRequest

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

	params := CreateCommunityParams{
		Name:        input.Name,
		OwnerID:     user.ID,
		Description: input.Description,
	}

	community, err := h.service.Create(r.Context(), params)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrDuplicateCommunityName):
			v.AddError("name", "community name is already in use")
			h.responder.ValidationError(w, v.Errors)
		default:
			h.responder.ServerError(w, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/api/v1/communities/%s", community.Name))

	view := communityToView(community)
	h.responder.JSON(w, http.StatusOK, toCommunityResponse(view))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("community_name")

	user, _ := request.GetUser(r)
	err := h.service.Delete(r.Context(), name, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			h.responder.NotFound(w, r)
		case errors.Is(err, errs.ErrForbidden):
			h.responder.Forbidden(w)
		default:
			h.responder.ServerError(w, err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
