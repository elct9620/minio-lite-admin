package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/elct9620/minio-lite-admin/internal/config"
	httpHandler "github.com/elct9620/minio-lite-admin/internal/handler/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed all:dist
var distFS embed.FS

func main() {
	// Load configuration
	cfg := config.Load()

	// Set up Chi router
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/health", httpHandler.HealthHandler)
		r.Get("/server-info", httpHandler.ServerInfoHandler)
	})

	// Frontend handler
	frontendHandler := httpHandler.NewFrontendHandler(cfg, distFS)
	r.Get("/*", frontendHandler.ServeHTTP)

	// Start server
	log.Printf("Server starting on %s", cfg.Server.Addr)
	if cfg.Server.Dev {
		log.Println("Running in development mode")
		log.Printf("Vite dev server URL: %s", cfg.Vite.URL)
		log.Println("Make sure to run 'npm run dev' for the Vite dev server")
	} else {
		log.Println("Running in production mode")
	}

	if err := http.ListenAndServe(cfg.Server.Addr, r); err != nil {
		log.Fatal(err)
	}
}
