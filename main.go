package main

import (
	"net/http"

	"github.com/elct9620/minio-lite-admin/internal/config"
	httpHandler "github.com/elct9620/minio-lite-admin/internal/handler/http"
	"github.com/elct9620/minio-lite-admin/internal/infra"
	"github.com/elct9620/minio-lite-admin/internal/logger"
	"github.com/elct9620/minio-lite-admin/internal/service"
)

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

	// Set up HTTP service with dependencies
	r, err := httpHandler.NewService(cfg, log, getServerInfoService, distFS)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create HTTP service")
	}

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
