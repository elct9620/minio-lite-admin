package main

import (
	"embed"
	"net/http"

	"github.com/elct9620/minio-lite-admin/internal/config"
	httpHandler "github.com/elct9620/minio-lite-admin/internal/handler/http"
	"github.com/elct9620/minio-lite-admin/internal/infra"
	"github.com/elct9620/minio-lite-admin/internal/logger"
	"github.com/elct9620/minio-lite-admin/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed all:dist
var distFS embed.FS

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	log := logger.New(logger.Config{
		Level:  cfg.Logger.Level,
		Pretty: cfg.Logger.Pretty,
	})
	logger.SetGlobalLogger(log)

	// Initialize MinIO client
	minioClient, err := infra.NewMinIOClient(infra.MinIOConfig{
		URL:      cfg.MinIO.URL,
		RootUser: cfg.MinIO.RootUser,
		Password: cfg.MinIO.Password,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize MinIO client")
	}

	// Initialize services
	getServerInfoService := service.NewGetServerInfoService(minioClient)

	// Set up Chi router
	r := chi.NewRouter()

	// Add middleware
	r.Use(httpHandler.Logger(log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/health", httpHandler.HealthHandler)
		r.Get("/server-info", httpHandler.ServerInfoHandler(getServerInfoService))
	})

	// Frontend handler
	frontendHandler := httpHandler.NewFrontendHandler(cfg, distFS)
	r.Get("/*", frontendHandler.ServeHTTP)

	// Start server
	log.Info().Str("addr", cfg.Server.Addr).Msg("Server starting")
	if cfg.Server.Dev {
		log.Info().Msg("Running in development mode")
		log.Info().Str("vite_url", cfg.Vite.URL).Msg("Vite dev server URL")
		log.Info().Msg("Make sure to run 'pnpm dev' for the Vite dev server")
	} else {
		log.Info().Msg("Running in production mode")
	}

	if err := http.ListenAndServe(cfg.Server.Addr, r); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
