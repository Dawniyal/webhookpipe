package database

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/Dawniyal/webhookpipe/internal/config"
	loggerConfig "github.com/Dawniyal/webhookpipe/internal/logger"
	pgxZero "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
)

const DatabasePingTimeout = 10

type Database struct {
	Pool *pgxpool.Pool
	Log  *zerolog.Logger
}

func New(cfg *config.Config, logger *zerolog.Logger) (*Database, error) {
	hostPort := net.JoinHostPort(cfg.Database.Host, strconv.Itoa(cfg.Database.Port))

	encodedPassword := url.QueryEscape(cfg.Database.Password)

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Database.User,
		encodedPassword,
		hostPort,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	pgxPoolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pgxPoolConfig.MaxConns = int32(cfg.Database.MaxOpenConns)
	pgxPoolConfig.MinConns = int32(cfg.Database.MaxIdleConns)
	pgxPoolConfig.MaxConnLifetime = cfg.Database.ConnMaxLifetime
	pgxPoolConfig.MaxConnIdleTime = cfg.Database.ConnMaxIdleTime

	globalLevel := logger.GetLevel()
	pgxLogger := loggerConfig.NewPgxLogger(globalLevel)

	pgxPoolConfig.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxZero.NewLogger(pgxLogger),
		LogLevel: tracelog.LogLevel(loggerConfig.GetPgxTraceLogLevel(globalLevel)),
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxPoolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), (DatabasePingTimeout * time.Second))
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info().Msg("connected to the database")

	db := &Database{
		Pool: pool,
		Log:  logger,
	}

	return db, nil
}
