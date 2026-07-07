package members

import (
	"errors"
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

func (h *Handler) Join(w http.ResponseWriter, r *http.Request) {
	user, _ := request.GetUser(r)
	community := request.GetCommunity(r)

	err := h.service.Join(r.Context(), community.ID, user.ID)
	if err != nil {
		switch {
		default:
			h.responder.ServerError(w, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Leave(w http.ResponseWriter, r *http.Request) {
	user, _ := request.GetUser(r)
	community := request.GetCommunity(r)

	err := h.service.Leave(r.Context(), community.ID, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrOwnershipTransferRequired):
			h.responder.ConflictMsg(w,
				"ownership_transfer_required",
				"transfer ownership before leaving the community",
			)
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
	pageParams := request.ParseCursorPagination(r, v, pagination.DecodeMemberCursor)

	if !v.Valid() {
		h.responder.ValidationError(w, v.Errors)
		return
	}

	limit := pageParams.Limit
	pageParams.Limit += 1

	params := ListParams{
		CommunityID: community.ID,
		Pagination:  pageParams,
	}

	memberships, err := h.service.List(r.Context(), params)
	if err != nil {
		h.responder.ServerError(w, err)
		return
	}

	var nextCursor string
	if len(memberships) > limit {
		page := memberships[:limit]
		last := page[len(page)-1]
		cursor := &pagination.MemberCursor{UserID: last.UserID, JoinedAt: last.JoinedAt}

		nextCursor = cursor.Encode()
		memberships = page
	}

	h.responder.JSON(w, http.StatusOK, toListMembersResponse(memberships, nextCursor))
}
