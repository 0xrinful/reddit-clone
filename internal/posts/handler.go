package posts

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
	id, err := request.ReadID(r)
	if err != nil {
		h.responder.NotFound(w, r)
		return
	}

	user, authenticated := request.GetUser(r)

	post, err := h.service.Get(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			h.responder.NotFound(w, r)
		default:
			h.responder.ServerError(w, err)
		}
		return
	}

	if authenticated && user.ID == post.UserID {
		h.responder.JSON(w, http.StatusOK, toPostOwnerResponse(post))
	} else {
		h.responder.JSON(w, http.StatusOK, toPostResponse(post))
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user, _ := request.GetUser(r)
	community := request.GetCommunity(r)

	var input CreatePostRequest

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
		Title:       input.Title,
		Body:        input.Body,
	}
	post, err := h.service.Create(r.Context(), params)
	if err != nil {
		h.responder.ServerError(w, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/api/v1/posts/%d", post.ID))

	view := postToView(post, community.Name)
	h.responder.JSON(w, http.StatusCreated, toPostOwnerResponse(view), headers)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := request.ReadID(r)
	if err != nil {
		h.responder.NotFound(w, r)
		return
	}

	var input UpdatePostRequest

	err = request.DecodeJSON(w, r, &input)
	if err != nil {
		h.responder.DecodeError(w, err)
		return
	}

	v := validator.New()
	if input.Validate(v); !v.Valid() {
		h.responder.ValidationError(w, v.Errors)
		return
	}

	user, _ := request.GetUser(r)
	params := UpdateParams(input)

	post, err := h.service.Update(r.Context(), id, user.ID, params)
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

	view := postToView(post, "")
	h.responder.JSON(w, http.StatusCreated, toPostOwnerResponse(view))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := request.ReadID(r)
	if err != nil {
		h.responder.NotFound(w, r)
		return
	}

	user, _ := request.GetUser(r)

	err = h.service.Delete(r.Context(), id, user.ID)
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

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	community := request.GetCommunity(r)

	v := validator.New()
	pageParams := request.ParsePagination(r, v, pagination.DecodePostCursor)

	sort := SortBy(request.ReadString(r.URL.Query(), "sort", "new"))
	v.Check(sort.IsValid(), "sort", "invalid sort value")

	if !v.Valid() {
		h.responder.ValidationError(w, v.Errors)
		return
	}

	limit := pageParams.Limit
	pageParams.Limit += 1 // used to determine if there is a next cursor

	params := ListParams{
		CommunityID: community.ID,
		Pagination:  pageParams,
		Sort:        sort,
	}

	posts, err := h.service.List(r.Context(), params)
	if err != nil {
		h.responder.ServerError(w, err)
		return
	}

	var nextCursor string
	if len(posts) > limit {
		page := posts[:limit] // trim extra row used for next page check
		last := page[len(page)-1]
		next := &pagination.PostCursor{ID: last.ID}

		switch sort {
		case SortByNew:
			next.CreatedAt = &last.CreatedAt
		case SortByTop, SortByHot:
			next.Score = &last.Score
		}

		nextCursor = next.Encode()
		posts = page
	}

	h.responder.JSON(w, http.StatusOK, toListPostsResponse(posts, nextCursor))
}
