package posts

import (
	"fmt"
	"net/http"

	"github.com/0xrinful/rush"

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

func (h *Handler) RegisterRoutes(r *rush.Router) {
	r.Post("/posts", h.Create)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID := 1 // for now
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

	post, err := h.service.CreatePost(r.Context(), int64(userID), community.ID, input)
	if err != nil {
		h.responder.ServerError(w, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/api/v1/r/%s/posts/%d", community.Name, post.ID))

	h.responder.JSON(w, http.StatusCreated, toPostResponse(post), headers)
}
