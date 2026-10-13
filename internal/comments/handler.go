package comments

import (
	"net/http"

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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user, _ := request.GetUser(r)

	postID, err := request.ReadPostID(r)
	if err != nil {
		h.responder.NotFound(w, r)
		return
	}

	var input CreateCommentRequest

	err = request.DecodeJSON(w, r, &input)
	if err != nil {
		h.responder.DecodeError(w, err)
		return
	}

	v := validator.New()
	if input.Vaildate(v); !v.Valid() {
		h.responder.ValidationError(w, v.Errors)
		return
	}

	params := CreateParams{
		UserID:   user.ID,
		PostID:   postID,
		Content:  input.Content,
		ParentID: nil,
	}

	comment, err := h.service.Create(r.Context(), params)
	if err != nil {
		h.responder.HandleServiceError(w, err)
		return
	}

	view := commentToView(comment)
	h.responder.JSON(w, http.StatusCreated, toCommentResponse(view))
}
