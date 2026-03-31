package middleware

import (
	"github.com/0xrinful/reddit-clone/internal/config"
	"github.com/0xrinful/reddit-clone/internal/shared/response"
)

type Middleware struct {
	responder *response.Responder
	config    config.Config
}

func New(responder *response.Responder, config config.Config) *Middleware {
	return &Middleware{responder: responder, config: config}
}
