package logger

import (
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Level  string `mapstructure:"level"`
	Pretty bool   `mapstructure:"pretty"`
}

func New(cfg Config) zerolog.Logger {
	level := parseLevel(cfg.Level)
	zerolog.SetGlobalLevel(level)

	var writer io.Writer = os.Stdout
	if cfg.Pretty {
		writer = zerolog.ConsoleWriter{Out: os.Stdout}
	}

	return zerolog.New(writer).With().
		Timestamp().
		Caller().
		Logger()
}

func parseLevel(levelStr string) zerolog.Level {
	switch strings.ToLower(levelStr) {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

func SetGlobalLogger(logger zerolog.Logger) {
	log.Logger = logger
}
