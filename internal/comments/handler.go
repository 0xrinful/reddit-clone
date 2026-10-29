package comments

import (
	"net/http"

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
		UserID:  user.ID,
		PostID:  postID,
		Content: input.Content,
	}

	comment, err := h.service.Create(r.Context(), params)
	if err != nil {
		h.responder.HandleServiceError(w, err)
		return
	}

	view := commentToView(comment)
	h.responder.JSON(w, http.StatusCreated, toCommentResponse(view))
}

func (h *Handler) Reply(w http.ResponseWriter, r *http.Request) {
	user, _ := request.GetUser(r)

	commentID, err := request.ReadCommentID(r)
	if err != nil {
		h.responder.NotFound(w, r)
		return
	}

	var input ReplyCommentRequest

	err = request.DecodeJSON(w, r, &input)
	if err != nil {
		h.responder.DecodeError(w, err)
		return
	}

	v := validator.New()
	if input.Vaildate(v); !v.Valid() {
		h.responder.DecodeError(w, err)
		return
	}

	params := ReplyParams{
		UserID:   user.ID,
		ParentID: commentID,
		Content:  input.Content,
	}

	comment, err := h.service.Reply(r.Context(), params)
	if err != nil {
		h.responder.HandleServiceError(w, err)
		return
	}

	view := commentToView(comment)
	h.responder.JSON(w, http.StatusCreated, toCommentResponse(view))
}

func (h *Handler) ListByPost(w http.ResponseWriter, r *http.Request) {
	postID, err := request.ReadPostID(r)
	if err != nil {
		h.responder.NotFound(w, r)
		return
	}

	v := validator.New()
	pageParams := request.ParseCursorPagination(r, v, pagination.DecodeCommentCursor)

	sort := SortBy(request.ReadString(r.URL.Query(), "sort", "top"))
	v.Check(sort.Valid(), "sort", "invalid sort value")

	if !v.Valid() {
		h.responder.ValidationError(w, v.Errors)
		return
	}

	limit := pageParams.Limit
	pageParams.Limit += 1

	params := ListByPostParams{
		Sort:       sort,
		Pagination: pageParams,
		PostID:     postID,
	}

	comments, err := h.service.ListByPost(r.Context(), params)
	if err != nil {
		h.responder.HandleServiceError(w, err)
		return
	}

	var nextCursor string
	if len(comments) > limit {
		page := comments[:limit]
		last := page[len(page)-1]
		next := &pagination.CommentCursor{ID: last.ID}

		switch sort {
		case SortByNew:
			next.CreatedAt = &last.CreatedAt
		case SortByTop:
			next.Score = &last.Score
		}

		nextCursor = next.Encode()
		comments = page
	}

	h.responder.JSON(w, http.StatusOK, toListCommentsResponse(comments, nextCursor))
}

func (h *Handler) ListReplies(w http.ResponseWriter, r *http.Request) {
	commendID, err := request.ReadCommentID(r)
	if err != nil {
		h.responder.NotFound(w, r)
		return
	}

	v := validator.New()
	pageParams := request.ParseCursorPagination(r, v, pagination.DecodeCommentCursor)

	if !v.Valid() {
		h.responder.ValidationError(w, v.Errors)
		return
	}

	limit := pageParams.Limit
	pageParams.Limit += 1

	params := ListRepliesParams{
		Pagination: pageParams,
		CommentID:  commendID,
	}

	comments, err := h.service.ListReplies(r.Context(), params)
	if err != nil {
		h.responder.HandleServiceError(w, err)
		return
	}

	var nextCursor string
	if len(comments) > limit {
		page := comments[:limit]
		last := page[len(page)-1]
		next := &pagination.CommentCursor{ID: last.ID, CreatedAt: &last.CreatedAt}
		nextCursor = next.Encode()
		comments = page
	}

	h.responder.JSON(w, http.StatusOK, toListCommentsResponse(comments, nextCursor))
}
