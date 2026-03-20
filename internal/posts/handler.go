package posts

import (
	"net/http"

	"github.com/0xrinful/rush"

	"github.com/0xrinful/reddit-clone/internal/shared/request"
	"github.com/0xrinful/reddit-clone/internal/shared/response"
)

type Handler struct {
	service Service
	respond *response.Responder
}

func NewHandler(svc Service, responder *response.Responder) *Handler {
	return &Handler{svc, responder}
}

func (h *Handler) RegisterRoutes(r *rush.Router) {
	r.Get("/posts/{id}", h.getPost)
	r.Post("/posts", h.createPost)
}

func (h *Handler) getPost(w http.ResponseWriter, r *http.Request) {
	post := &Post{
		ID:    10,
		Title: "hehe",
		Body:  "hehe",
	}
	h.respond.JSON(w, http.StatusOK, toPostResponse(post))
}

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	var input CreatePostRequest

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		h.respond.DecodeError(w, err)
		return
	}

	post := &Post{
		Title: input.Title,
		Body:  input.Body,
	}

	h.respond.JSON(w, http.StatusOK, toPostResponse(post))
}
