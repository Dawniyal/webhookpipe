package middleware

import (
	"context"
	"net/http"

	"github.com/rs/zerolog"
)

const (
	LoggerKey = "logger"
)

type ContextEnhancer struct {
	logger *zerolog.Logger
}

func NewContextEnhancer(logger *zerolog.Logger) *ContextEnhancer {
	return &ContextEnhancer{logger: logger}
}

func (ce *ContextEnhancer) EnhanceContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := GetRequestID(r.Context())

		contextLogger := ce.logger.With().
			Str("request_id", requestID).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Logger()

		ctx := context.WithValue(r.Context(), LoggerKey, &contextLogger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetLogger(reqCtx context.Context) *zerolog.Logger {
	if logger := reqCtx.Value(LoggerKey).(*zerolog.Logger); logger != nil {
		return logger
	}

	logger := zerolog.Nop()
	return &logger
}
