package config

import (
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	Server        ServerConfig        `koanf:"server" validate:"required"`
	Database      DatabaseConfig      `koanf:"database" validate:"required"`
	Observability ObservabilityConfig `koanf:"observability" validate:"required"`
}

type ServerConfig struct {
	Port               string        `koanf:"port" validate:"required"`
	ReadTimeout        time.Duration `koanf:"read_timeout" validate:"required"`
	WriteTimeout       time.Duration `koanf:"write_timeout" validate:"required"`
	IdleTimeout        time.Duration `koanf:"idle_timeout" validate:"required"`
	CORSAllowedOrigins []string      `koanf:"cors_allowed_origins" validate:"required"`
}

type DatabaseConfig struct {
	Host            string        `koanf:"host" validate:"required"`
	Port            int           `koanf:"port" validate:"required"`
	User            string        `koanf:"user" validate:"required"`
	Password        string        `koanf:"password"`
	Name            string        `koanf:"name" validate:"required"`
	SSLMode         string        `koanf:"ssl_mode" validate:"required"`
	MaxOpenConns    int           `koanf:"max_open_conns" validate:"required"`
	MaxIdleConns    int           `koanf:"max_idle_conns" validate:"required"`
	ConnMaxLifetime time.Duration `koanf:"conn_max_lifetime" validate:"required"`
	ConnMaxIdleTime time.Duration `koanf:"conn_max_idle_time" validate:"required"`
}

type ObservabilityConfig struct {
	Logging      LoggingConfig      `koanf:"logging" validate:"required"`
	HealthChecks HealthChecksConfig `koanf:"health_checks" validate:"required"`
}

type LoggingConfig struct {
	Level  string `koanf:"level" validate:"required,oneof=debug info warn error"`
	Format string `koanf:"format" validate:"required,oneof=json console"`
}

type HealthChecksConfig struct {
	Enabled  bool          `koanf:"enabled" validate:"required"`
	Interval time.Duration `koanf:"interval" validate:"required,min=1s"`
	Timeout  time.Duration `koanf:"timeout" validate:"required,min=1s"`
}

func LoadConfig(path string) (*Config, error) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	err := godotenv.Load(path)
	if err != nil {
		logger.Fatal().Err(err)
	}

	k := koanf.New(".")

	err = k.Load(env.Provider("WEBHOOKPIPE_", ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "WEBHOOKPIPE_"))
	}), nil)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not load initial env variables")
	}

	cfg := &Config{}

	err = k.Unmarshal("", cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not unmarshal main config")
	}

	validate := validator.New()

	err = validate.Struct(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("config validation failed")
	}

	return cfg, nil
}
