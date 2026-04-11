package middleware

import (
	"sync"

	"github.com/0xrinful/reddit-clone/internal/config"
	"github.com/0xrinful/reddit-clone/internal/shared/response"
)

type Middleware struct {
	responder *response.Responder
	config    config.Config
	mu        sync.RWMutex
	clients   map[string]*client
}

func New(responder *response.Responder, config config.Config) *Middleware {
	m := &Middleware{
		responder: responder,
		config:    config,
		clients:   make(map[string]*client),
	}
	go m.startClientCleanup()
	return m
}
