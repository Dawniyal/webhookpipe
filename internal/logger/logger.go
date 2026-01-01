package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Dawniyal/webhookpipe/internal/config"

	"github.com/rs/zerolog"
)

func NewLogger(cfg *config.Config) zerolog.Logger {
	var logLevel zerolog.Level

	switch cfg.Observability.Logging.Level {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	default:
		logLevel = zerolog.InfoLevel
	}

	zerolog.TimeFieldFormat = time.RFC3339

	var writer io.Writer

	if cfg.Observability.Logging.Format == "json" {
		writer = os.Stdout
	} else {
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		writer = consoleWriter
	}

	logger := zerolog.New(writer).
		Level(logLevel).
		With().
		Timestamp().
		Logger()

	return logger
}

func NewPgxLogger(level zerolog.Level) zerolog.Logger {
	writer := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		FormatFieldValue: func(i any) string {
			switch v := i.(type) {
			case string:
				if len(v) > 200 {
					return v[:200] + "..."
				}
				return v
			case []byte:
				var obj any
				if err := json.Unmarshal(v, &obj); err == nil {
					if pretty, err := json.MarshalIndent(obj, "", "\t"); err == nil {
						return "\n" + string(pretty)
					} else {
						return err.Error()
					}

				}
				return string(v)
			default:
				return fmt.Sprintf("%v", v)

			}
		},
	}

	return zerolog.New(writer).
		With().
		Timestamp().
		Str("component", "database").Logger()
}

// GetPgxTraceLogLevel converts zerolog level to pgx tracelog level
func GetPgxTraceLogLevel(level zerolog.Level) int {
	switch level {
	case zerolog.DebugLevel:
		return 6 // tracelog.LogLevelDebug
	case zerolog.InfoLevel:
		return 4 // tracelog.LogLevelInfo
	case zerolog.WarnLevel:
		return 3 // tracelog.LogLevelWarn
	case zerolog.ErrorLevel:
		return 2 // tracelog.LogLevelError
	default:
		return 0 // tracelog.LogLevelNone
	}
}
