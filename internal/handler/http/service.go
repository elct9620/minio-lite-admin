package http

import (
	"embed"
	"net/http"

	"github.com/elct9620/minio-lite-admin/internal/config"
	"github.com/elct9620/minio-lite-admin/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

// Service handles all HTTP requests and contains all dependencies
type Service struct {
	config               *config.Config
	logger               zerolog.Logger
	getServerInfoService *service.GetServerInfoService
	distFS               embed.FS
}

// NewService creates a new HTTP service with all dependencies and returns the configured router
func NewService(
	cfg *config.Config,
	logger zerolog.Logger,
	getServerInfoService *service.GetServerInfoService,
	distFS embed.FS,
) (http.Handler, error) {
	svc := &Service{
		config:               cfg,
		logger:               logger,
		getServerInfoService: getServerInfoService,
		distFS:               distFS,
	}

	router := chi.NewRouter()

	// Add middleware
	router.Use(Logger(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)

	// API routes
	router.Route("/api", func(r chi.Router) {
		r.Get("/health", svc.GetHealthHandler)
		r.Get("/server-info", svc.GetServerInfoHandler)
	})

	// Frontend routes
	router.Get("/*", svc.GetRootHandler)

	return router, nil
}
