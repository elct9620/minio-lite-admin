package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func Logger(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{logger})
}

type StructuredLogger struct {
	Logger zerolog.Logger
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{Logger: l.Logger}
	logFields := zerolog.Dict().
		Str("method", r.Method).
		Str("uri", r.RequestURI).
		Str("remote_addr", r.RemoteAddr).
		Str("user_agent", r.UserAgent())

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields = logFields.Str("req_id", reqID)
	}

	entry.Logger = l.Logger.With().Dict("request", logFields).Logger()
	entry.Logger.Info().Msg("request started")

	return entry
}

type StructuredLoggerEntry struct {
	Logger zerolog.Logger
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger.Info().
		Int("status", status).
		Int("bytes", bytes).
		Dur("elapsed", elapsed).
		Msg("request completed")
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger.Error().
		Str("stack", string(stack)).
		Interface("panic", v).
		Msg("request panicked")
}
