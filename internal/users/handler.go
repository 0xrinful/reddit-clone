package users

import (
	"github.com/0xrinful/reddit-clone/internal/shared/response"
)

type Handler struct {
	service   Service
	responder *response.Responder
}

func NewHandler(svc Service, responder *response.Responder) *Handler {
	return &Handler{svc, responder}
}
