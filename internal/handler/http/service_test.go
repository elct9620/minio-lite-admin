package http

import (
	"embed"
	"os"

	"github.com/elct9620/minio-lite-admin/internal/config"
	"github.com/elct9620/minio-lite-admin/internal/service"
	"github.com/elct9620/minio-lite-admin/internal/testability/minio"
	"github.com/rs/zerolog"
)

// testService creates a Service instance for testing purposes
func testService() *Service {
	// Create test logger (discard output to avoid noise in tests)
	logger := zerolog.New(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.Out = os.Stdout
	})).Level(zerolog.Disabled) // Disable logging in tests

	// Create minimal test config
	cfg := &config.Config{
		Server: config.Server{
			Addr: ":8080",
			Dev:  true,
		},
		Vite: config.Vite{
			URL:   "http://localhost:5173",
			Entry: "/src/main.ts",
		},
	}

	// Create mock service (for now, nil since health endpoint doesn't use it)
	var getServerInfoService *service.GetServerInfoService

	// Create empty embed.FS for testing
	var distFS embed.FS

	return &Service{
		config:               cfg,
		logger:               logger,
		getServerInfoService: getServerInfoService,
		distFS:               distFS,
	}
}

// testServiceWithMockMinIO creates a Service instance with mock MinIO server for testing
func testServiceWithMockMinIO() (*Service, *minio.MockMinIOServer) {
	// Create test logger (discard output to avoid noise in tests)
	logger := zerolog.New(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.Out = os.Stdout
	})).Level(zerolog.Disabled) // Disable logging in tests

	// Create minimal test config
	cfg := &config.Config{
		Server: config.Server{
			Addr: ":8080",
			Dev:  true,
		},
		Vite: config.Vite{
			URL:   "http://localhost:5173",
			Entry: "/src/main.ts",
		},
	}

	// Create mock MinIO server
	mockMinIO := minio.NewMockMinIOServer()

	// Create MinIO admin client pointing to mock server
	minioClient, err := mockMinIO.CreateMinIOClient()
	if err != nil {
		panic("Failed to create mock MinIO client: " + err.Error())
	}

	// Create service with mock MinIO client
	getServerInfoService := service.NewGetServerInfoService(minioClient)

	// Create empty embed.FS for testing
	var distFS embed.FS

	svc := &Service{
		config:               cfg,
		logger:               logger,
		getServerInfoService: getServerInfoService,
		distFS:               distFS,
	}

	return svc, mockMinIO
}