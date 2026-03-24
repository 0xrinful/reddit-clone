package middleware

import "github.com/0xrinful/reddit-clone/internal/shared/response"

type Middleware struct {
	responder *response.Responder
}

func New(responder *response.Responder) *Middleware {
	return &Middleware{responder: responder}
}
