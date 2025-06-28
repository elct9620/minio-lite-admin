package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/elct9620/minio-lite-admin/internal/config"
	"github.com/elct9620/minio-lite-admin/internal/service"
	"github.com/elct9620/minio-lite-admin/internal/testability/minio"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

func TestService_DeleteAccessKeysHandler(t *testing.T) {
	tests := []struct {
		name               string
		setupMock          func(*minio.MockMinIOServer)
		accessKey          string
		expectedStatusCode int
		expectedError      string
		validateResponse   func(t *testing.T, response *service.DeleteServiceAccountResponse)
	}{
		{
			name: "successful deletion",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetDeleteServiceAccountSuccess()
			},
			accessKey:          "AKIAIOSFODNN7EXAMPLE",
			expectedStatusCode: http.StatusOK,
			validateResponse: func(t *testing.T, response *service.DeleteServiceAccountResponse) {
				if response.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", response.AccessKey)
				}
				if response.Message != "Service account deleted successfully" {
					t.Errorf("Expected Message %q, got %q", "Service account deleted successfully", response.Message)
				}
			},
		},
		{
			name: "successful deletion with different access key",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetDeleteServiceAccountSuccess()
			},
			accessKey:          "AKIAI44QH8DHBEXAMPLE",
			expectedStatusCode: http.StatusOK,
			validateResponse: func(t *testing.T, response *service.DeleteServiceAccountResponse) {
				if response.AccessKey != "AKIAI44QH8DHBEXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAI44QH8DHBEXAMPLE", response.AccessKey)
				}
				if response.Message != "Service account deleted successfully" {
					t.Errorf("Expected Message %q, got %q", "Service account deleted successfully", response.Message)
				}
			},
		},
		{
			name:               "missing access key parameter",
			setupMock:          func(mock *minio.MockMinIOServer) {},
			accessKey:          "",
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      "Access key is required",
		},
		{
			name:               "empty access key parameter",
			setupMock:          func(mock *minio.MockMinIOServer) {},
			accessKey:          "   ",
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      "Access key is required",
		},
		{
			name: "MinIO server error",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetDeleteServiceAccountError(500, "Internal Server Error")
			},
			accessKey:          "AKIAIOSFODNN7EXAMPLE",
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      "Failed to delete access key",
		},
		{
			name: "MinIO authentication error",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetDeleteServiceAccountError(401, "Unauthorized")
			},
			accessKey:          "AKIAIOSFODNN7EXAMPLE",
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      "Failed to delete access key",
		},
		{
			name: "MinIO not found error",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetDeleteServiceAccountError(404, "Service account not found")
			},
			accessKey:          "NONEXISTENT",
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      "Failed to delete access key",
		},
		{
			name: "MinIO forbidden error",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetDeleteServiceAccountError(403, "Forbidden")
			},
			accessKey:          "AKIAIOSFODNN7EXAMPLE",
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      "Failed to delete access key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock MinIO server
			mockServer := minio.NewMockMinIOServer()
			defer mockServer.Close()

			tt.setupMock(mockServer)

			// Create MinIO client
			minioClient, err := mockServer.CreateMinIOClient()
			if err != nil {
				t.Fatalf("Failed to create MinIO client: %v", err)
			}

			// Create service dependencies
			deleteServiceAccountService := service.NewDeleteServiceAccountService(minioClient)

			// Create HTTP service
			testService := createTestServiceForDelete(t, deleteServiceAccountService)

			// Create request with Chi context for URL parameters
			// Use a safe URL path - the actual parameter comes from Chi context
			req := httptest.NewRequest(http.MethodDelete, "/api/access-keys/test", nil)

			// Set up Chi context with URL parameters
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("accessKey", tt.accessKey)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Add context with logger
			ctx := zerolog.New(zerolog.NewTestWriter(t)).WithContext(req.Context())
			req = req.WithContext(ctx)

			// Create response recorder
			w := httptest.NewRecorder()

			// Execute request
			testService.DeleteAccessKeysHandler(w, req)

			// Validate status code
			if w.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, w.Code)
			}

			if tt.expectedError != "" {
				// Validate error response
				if !strings.Contains(w.Body.String(), tt.expectedError) {
					t.Errorf("Expected response to contain %q, got %q", tt.expectedError, w.Body.String())
				}
			} else {
				// Validate success response
				contentType := w.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Expected Content-Type %q, got %q", "application/json", contentType)
				}

				var response service.DeleteServiceAccountResponse
				err := json.NewDecoder(w.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if tt.validateResponse != nil {
					tt.validateResponse(t, &response)
				}
			}
		})
	}
}

func TestService_DeleteAccessKeysHandler_UrlParameterEdgeCases(t *testing.T) {
	tests := []struct {
		name               string
		setupRequest       func() *http.Request
		expectedStatusCode int
		expectedError      string
	}{
		{
			name: "URL parameter with special characters",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodDelete, "/api/access-keys/test", nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("accessKey", "AKIA TEST") // URL decoded
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
				return req
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "URL parameter with only whitespace",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodDelete, "/api/access-keys/test", nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("accessKey", "   ")
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
				return req
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      "Access key is required",
		},
		{
			name: "URL parameter with leading/trailing spaces",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodDelete, "/api/access-keys/test", nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("accessKey", "  AKIATEST  ")
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
				return req
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock MinIO server
			mockServer := minio.NewMockMinIOServer()
			defer mockServer.Close()

			mockServer.SetDeleteServiceAccountSuccess()

			// Create MinIO client
			minioClient, err := mockServer.CreateMinIOClient()
			if err != nil {
				t.Fatalf("Failed to create MinIO client: %v", err)
			}

			// Create service dependencies
			deleteServiceAccountService := service.NewDeleteServiceAccountService(minioClient)

			// Create HTTP service
			testService := createTestServiceForDelete(t, deleteServiceAccountService)

			// Setup request
			req := tt.setupRequest()

			// Add context with logger
			ctx := zerolog.New(zerolog.NewTestWriter(t)).WithContext(req.Context())
			req = req.WithContext(ctx)

			// Create response recorder
			w := httptest.NewRecorder()

			// Execute request
			testService.DeleteAccessKeysHandler(w, req)

			// Validate status code
			if w.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, w.Code)
			}

			if tt.expectedError != "" {
				if !strings.Contains(w.Body.String(), tt.expectedError) {
					t.Errorf("Expected response to contain %q, got %q", tt.expectedError, w.Body.String())
				}
			}
		})
	}
}

// createTestServiceForDelete creates a Service instance for testing delete functionality
func createTestServiceForDelete(
	t *testing.T,
	deleteServiceAccountService *service.DeleteServiceAccountService,
) *Service {
	// Create test config
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

	// Create test logger
	logger := zerolog.New(zerolog.NewTestWriter(t))

	// Create service with test dependencies
	return &Service{
		config:                      cfg,
		logger:                      logger,
		deleteServiceAccountService: deleteServiceAccountService,
	}
}

// TestService_DeleteAccessKeysHandler_Integration tests the handler through the full router
func TestService_DeleteAccessKeysHandler_Integration(t *testing.T) {
	// Setup mock MinIO server
	mockServer := minio.NewMockMinIOServer()
	defer mockServer.Close()

	mockServer.SetDeleteServiceAccountSuccess()

	// Create MinIO client
	minioClient, err := mockServer.CreateMinIOClient()
	if err != nil {
		t.Fatalf("Failed to create MinIO client: %v", err)
	}

	// Create service dependencies
	deleteServiceAccountService := service.NewDeleteServiceAccountService(minioClient)

	// Create full service
	testService := createTestServiceForDelete(t, deleteServiceAccountService)

	// Create router and add our handler
	router := chi.NewRouter()
	router.Delete("/api/access-keys/{accessKey}", testService.DeleteAccessKeysHandler)

	// Create test server
	server := httptest.NewServer(router)
	defer server.Close()

	// Make request
	req, err := http.NewRequest(http.MethodDelete, server.URL+"/api/access-keys/AKIAIOSFODNN7EXAMPLE", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Validate response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type %q, got %q", "application/json", contentType)
	}

	var response service.DeleteServiceAccountResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
		t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", response.AccessKey)
	}
	if response.Message != "Service account deleted successfully" {
		t.Errorf("Expected Message %q, got %q", "Service account deleted successfully", response.Message)
	}
}
