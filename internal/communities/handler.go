package communities

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/0xrinful/reddit-clone/internal/shared/errs"
	"github.com/0xrinful/reddit-clone/internal/shared/pagination"
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

	params := CreateParams{
		Name:        input.Name,
		OwnerID:     user.ID,
		Description: input.Description,
	}

	community, err := h.service.Create(r.Context(), params)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrDuplicateCommunityName):
			v.AddError("name", "community already exists")
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
		case errors.Is(err, errs.ErrPermissionDenied):
			h.responder.PermissionDenied(w)
		default:
			h.responder.ServerError(w, err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("community_name")
	user, _ := request.GetUser(r)

	var input UpdateCommunitytRequest
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

	p := UpdateParams(input)
	community, err := h.service.Update(r.Context(), name, user.ID, p)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			h.responder.NotFound(w, r)
		case errors.Is(err, errs.ErrPermissionDenied):
			h.responder.PermissionDenied(w)
		case errors.Is(err, errs.ErrDuplicateCommunityName):
			v.AddError("name", "community name is already in use")
			h.responder.ValidationError(w, v.Errors)
		default:
			h.responder.ServerError(w, err)
		}
		return
	}

	view := communityToView(community)
	h.responder.JSON(w, http.StatusOK, toCommunityResponse(view))
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	q := request.ReadString(r.URL.Query(), "q", "")

	v := validator.New()
	pageParams := request.ParseOffsetPagination(r, v)

	if !v.Valid() {
		h.responder.ValidationError(w, v.Errors)
		return
	}

	params := SearchParams{
		Name:       q,
		Pagination: pageParams,
	}

	communities, err := h.service.Search(r.Context(), params)
	if err != nil {
		h.responder.ServerError(w, err)
		return
	}

	h.responder.JSON(w, http.StatusOK, toSearchCommunitiesResponse(communities))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	v := validator.New()
	pageParams := request.ParseCursorPagination(r, v, pagination.DecodeCommunityCursor)

	if !v.Valid() {
		h.responder.ValidationError(w, v.Errors)
		return
	}

	limit := pageParams.Limit
	pageParams.Limit += 1 // used to determine if there is a next cursor

	params := ListParams{
		Pagination: pageParams,
	}

	communities, err := h.service.List(r.Context(), params)
	if err != nil {
		h.responder.ServerError(w, err)
		return
	}

	var nextCursor string
	if len(communities) > limit {
		page := communities[:limit] // trim extra row used for next page check
		last := page[len(page)-1]
		next := &pagination.CommunityCursor{ID: last.ID, CreatedAt: &last.CreatedAt}

		nextCursor = next.Encode()
		communities = page
	}

	h.responder.JSON(w, http.StatusOK, toListCommunitiesResponse(communities, nextCursor))
}
