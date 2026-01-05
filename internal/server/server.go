package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Dawniyal/webhookpipe/internal/config"
	"github.com/Dawniyal/webhookpipe/internal/database"
	"github.com/Dawniyal/webhookpipe/internal/lib/job"
	"github.com/Dawniyal/webhookpipe/internal/server"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type Server struct {
	Config     *config.Config
	Logger     *zerolog.Logger
	DB         *database.Database
	Redis      *redis.Client
	Job        *job.JobService
	httpServer *http.Server
}

func New(cfg *config.Config, logger *zerolog.Logger) (*Server, error) {
	db, err := database.New(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	redis := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redis.Ping(ctx).Err(); err != nil {
		logger.Error().Err(err).Msg("Failed to connect to Redis")
		return nil, err
	}

	server := &Server{
		Config: cfg,
		Logger: logger,
		DB:     db,
		Redis:  redis,
		Job:    jobService,
	}

	return server, nil
}

func (s *Server) SetupHttpServer(handler http.Handler) {
	s.httpServer = &http.Server{
		Addr:         ":" + s.Config.Server.Port,
		Handler:      handler,
		ReadTimeout:  s.Config.Server.ReadTimeout,
		WriteTimeout: s.Config.Server.WriteTimeout,
		IdleTimeout:  s.Config.Server.IdleTimeout,
	}
}

func (s *Server) Start() error {
	if s.httpServer == nil {
		return errors.New("HTTP server not initialized")
	}

	s.Logger.Info().
		Str("port", s.Config.Server.Port).
		Msg("starting server")

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	if err := s.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	if s.Job != nil {
		s.Job.Stop()
	}

	return nil
}
