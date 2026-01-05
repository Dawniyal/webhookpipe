package middleware

import (
	"github.com/Dawniyal/webhookpipe/internal/server"
)

type Middlewares struct {
	ContextEnhancer *ContextEnhancer
}

func NewMiddlewares(s *server.Server) *Middlewares {
	return &Middlewares{
		ContextEnhancer: NewContextEnhancer(s.Logger),
	}
}
