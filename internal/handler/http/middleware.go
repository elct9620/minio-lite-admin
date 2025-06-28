package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func Logger() func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{})
}

type StructuredLogger struct{}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	// Use context-aware logger that already includes request ID
	contextLogger := zerolog.Ctx(r.Context())

	// Create request-specific logger with expanded fields
	entry := &StructuredLoggerEntry{
		Logger: contextLogger.With().
			Str("method", r.Method).
			Str("uri", r.RequestURI).
			Str("remote_addr", r.RemoteAddr).
			Str("user_agent", r.UserAgent()).
			Logger(),
	}

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
		Send()
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger.Error().
		Str("stack", string(stack)).
		Interface("panic", v).
		Send()
}

// ContextLogger middleware adds the logger to the request context for use in handlers
func ContextLogger(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Create a logger with request ID if available
			contextLogger := logger
			if reqID := middleware.GetReqID(ctx); reqID != "" {
				contextLogger = logger.With().Str("req_id", reqID).Logger()
			}

			// Add logger to context
			ctx = contextLogger.WithContext(ctx)

			// Continue with the enhanced context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
