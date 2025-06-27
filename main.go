package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Create HTTP server
	server := &http.Server{
		Addr:    cfg.Server.Addr,
		Handler: r,
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

	// Channel to listen for interrupt signal to terminate server gracefully
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	log.Info().Msg("Server started successfully")

	// Wait for interrupt signal
	<-quit
	log.Info().Msg("Shutting down server...")

	// Create a deadline for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
		return
	}

	log.Info().Msg("Server shutdown complete")
}
