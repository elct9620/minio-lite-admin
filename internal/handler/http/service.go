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
	config                      *config.Config
	logger                      zerolog.Logger
	getServerInfoService        *service.GetServerInfoService
	listAccessKeysService       *service.ListAccessKeysService
	addServiceAccountService    *service.AddServiceAccountService
	deleteServiceAccountService *service.DeleteServiceAccountService
	updateServiceAccountService *service.UpdateServiceAccountService
	distFS                      embed.FS
}

// NewService creates a new HTTP service with all dependencies and returns the configured router
func NewService(
	cfg *config.Config,
	logger zerolog.Logger,
	getServerInfoService *service.GetServerInfoService,
	listAccessKeysService *service.ListAccessKeysService,
	addServiceAccountService *service.AddServiceAccountService,
	deleteServiceAccountService *service.DeleteServiceAccountService,
	updateServiceAccountService *service.UpdateServiceAccountService,
	distFS embed.FS,
) (http.Handler, error) {
	svc := &Service{
		config:                      cfg,
		logger:                      logger,
		getServerInfoService:        getServerInfoService,
		listAccessKeysService:       listAccessKeysService,
		addServiceAccountService:    addServiceAccountService,
		deleteServiceAccountService: deleteServiceAccountService,
		updateServiceAccountService: updateServiceAccountService,
		distFS:                      distFS,
	}

	router := chi.NewRouter()

	// Add middleware
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(ContextLogger(logger))
	router.Use(Logger())

	// API routes
	router.Route("/api", func(r chi.Router) {
		r.Get("/health", svc.GetHealthHandler)
		r.Get("/server-info", svc.GetServerInfoHandler)
		r.Get("/data-usage", svc.GetDataUsageHandler)
		r.Get("/access-keys", svc.GetAccessKeysHandler)
		r.Post("/access-keys", svc.PostAccessKeysHandler)
		r.Put("/access-keys/{accessKey}", svc.PutAccessKeysHandler)
		r.Delete("/access-keys/{accessKey}", svc.DeleteAccessKeysHandler)
	})

	// Frontend routes
	router.Get("/*", svc.GetRootHandler)

	return router, nil
}
