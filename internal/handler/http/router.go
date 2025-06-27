package http

import (
	"embed"

	"github.com/elct9620/minio-lite-admin/internal/config"
	"github.com/elct9620/minio-lite-admin/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

// RouterDependencies holds all dependencies needed to create the router
type RouterDependencies struct {
	Config               *config.Config
	Logger               zerolog.Logger
	GetServerInfoService *service.GetServerInfoService
	DistFS               embed.FS
}

// NewRouter creates and configures the main HTTP router with all routes and middleware
func NewRouter(deps RouterDependencies) *chi.Mux {
	r := chi.NewRouter()

	// Add middleware
	r.Use(Logger(deps.Logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/health", HealthHandler)
		r.Get("/server-info", ServerInfoHandler(deps.GetServerInfoService))
	})

	// Frontend handler
	frontendHandler := NewFrontendHandler(deps.Config, deps.DistFS)
	r.Get("/*", frontendHandler.ServeHTTP)

	return r
}