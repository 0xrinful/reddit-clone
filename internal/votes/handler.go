package votes

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

func (h *Handler) VotePost(w http.ResponseWriter, r *http.Request) {
	user, _ := request.GetUser(r)

	postID, err := request.ReadPostID(r)
	if err != nil {
		h.responder.NotFound(w, r)
		return
	}

	var input VotePostRequest
	if err := request.DecodeJSON(w, r, &input); err != nil {
		h.responder.DecodeError(w, err)
		return
	}

	v := validator.New()
	if input.Validate(v); !v.Valid() {
		h.responder.ValidationError(w, v.Errors)
		return
	}

	vote := PostVote{
		UserID: user.ID,
		PostID: postID,
		Value:  VoteValue(input.Value),
	}

	result, err := h.service.VotePost(r.Context(), vote)
	if err != nil {
		h.responder.HandleServiceError(w, err)
		return
	}

	h.responder.JSON(w, http.StatusOK, toVotePostResponse(result))
}

func (h *Handler) UnvotePost(w http.ResponseWriter, r *http.Request) {
	user, _ := request.GetUser(r)

	postID, err := request.ReadPostID(r)
	if err != nil {
		h.responder.NotFound(w, r)
		return
	}

	result, err := h.service.UnvotePost(r.Context(), user.ID, postID)
	if err != nil {
		h.responder.HandleServiceError(w, err)
		return
	}

	h.responder.JSON(w, http.StatusOK, toVotePostResponse(result))
}
