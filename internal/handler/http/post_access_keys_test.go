package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/elct9620/minio-lite-admin/internal/config"
	"github.com/elct9620/minio-lite-admin/internal/service"
	"github.com/elct9620/minio-lite-admin/internal/testability/minio"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

func TestService_PostAccessKeysHandler(t *testing.T) {
	tests := []struct {
		name               string
		setupMock          func(*minio.MockMinIOServer)
		requestBody        interface{}
		expectedStatusCode int
		expectedError      string
		validateResponse   func(t *testing.T, response *service.CreateServiceAccountResponse)
	}{
		{
			name: "successful creation with auto-generated keys",
			setupMock: func(mock *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				mock.SetAddServiceAccountResponse(scenarios.SuccessfulAddServiceAccount())
			},
			requestBody: service.CreateServiceAccountRequest{
				Name:        "test-service-account",
				Description: "Test service account",
			},
			expectedStatusCode: http.StatusCreated,
			validateResponse: func(t *testing.T, response *service.CreateServiceAccountResponse) {
				if response.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", response.AccessKey)
				}
				if response.SecretKey != "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" {
					t.Errorf("Expected SecretKey %q, got %q", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", response.SecretKey)
				}
				if response.Name != "test-service-account" {
					t.Errorf("Expected Name %q, got %q", "test-service-account", response.Name)
				}
				if response.Description != "Test service account" {
					t.Errorf("Expected Description %q, got %q", "Test service account", response.Description)
				}
			},
		},
		{
			name: "successful creation without name (MinIO auto-generates)",
			setupMock: func(mock *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				mock.SetAddServiceAccountResponse(scenarios.SuccessfulAddServiceAccount())
			},
			requestBody: service.CreateServiceAccountRequest{
				Description: "Service account without explicit name",
			},
			expectedStatusCode: http.StatusCreated,
			validateResponse: func(t *testing.T, response *service.CreateServiceAccountResponse) {
				if response.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", response.AccessKey)
				}
				if response.SecretKey != "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" {
					t.Errorf("Expected SecretKey %q, got %q", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", response.SecretKey)
				}
				// Name should be empty when not provided
				if response.Name != "" {
					t.Errorf("Expected empty Name, got %q", response.Name)
				}
				if response.Description != "Service account without explicit name" {
					t.Errorf("Expected Description %q, got %q", "Service account without explicit name", response.Description)
				}
			},
		},
		{
			name: "successful creation with custom keys",
			setupMock: func(mock *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				mock.SetAddServiceAccountResponse(scenarios.CustomKeysAddServiceAccount("CUSTOM123", "customSecret456"))
			},
			requestBody: service.CreateServiceAccountRequest{
				Name:        "custom-service-account",
				Description: "Custom service account",
				AccessKey:   "CUSTOM123",
				SecretKey:   "customSecret456",
			},
			expectedStatusCode: http.StatusCreated,
			validateResponse: func(t *testing.T, response *service.CreateServiceAccountResponse) {
				if response.AccessKey != "CUSTOM123" {
					t.Errorf("Expected AccessKey %q, got %q", "CUSTOM123", response.AccessKey)
				}
				if response.SecretKey != "customSecret456" {
					t.Errorf("Expected SecretKey %q, got %q", "customSecret456", response.SecretKey)
				}
				if response.Name != "custom-service-account" {
					t.Errorf("Expected Name %q, got %q", "custom-service-account", response.Name)
				}
				if response.Description != "Custom service account" {
					t.Errorf("Expected Description %q, got %q", "Custom service account", response.Description)
				}
			},
		},
		{
			name: "successful creation with policy",
			setupMock: func(mock *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				mock.SetAddServiceAccountResponse(scenarios.SuccessfulAddServiceAccount())
			},
			requestBody: service.CreateServiceAccountRequest{
				Name:   "policy-service-account",
				Policy: `{"Version": "2012-10-17", "Statement": [{"Effect": "Allow", "Action": ["s3:GetObject"], "Resource": ["arn:aws:s3:::bucket/*"]}]}`,
			},
			expectedStatusCode: http.StatusCreated,
			validateResponse: func(t *testing.T, response *service.CreateServiceAccountResponse) {
				if response.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", response.AccessKey)
				}
				if response.Name != "policy-service-account" {
					t.Errorf("Expected Name %q, got %q", "policy-service-account", response.Name)
				}
			},
		},
		{
			name: "successful creation with expiration",
			setupMock: func(mock *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				expiration := time.Now().Add(24 * time.Hour)
				mock.SetAddServiceAccountResponse(scenarios.ExpiringAddServiceAccount(expiration))
			},
			requestBody: service.CreateServiceAccountRequest{
				Name: "expiring-service-account",
				Expiration: func() *int64 {
					exp := time.Now().Add(24 * time.Hour).Unix()
					return &exp
				}(),
			},
			expectedStatusCode: http.StatusCreated,
			validateResponse: func(t *testing.T, response *service.CreateServiceAccountResponse) {
				if response.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", response.AccessKey)
				}
				if response.Name != "expiring-service-account" {
					t.Errorf("Expected Name %q, got %q", "expiring-service-account", response.Name)
				}
				if response.Expiration.IsZero() {
					t.Errorf("Expected non-zero Expiration, got zero")
				}
			},
		},
		{
			name: "missing name handled by MinIO",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetAddServiceAccountError(400, "The service account name is invalid")
			},
			requestBody:        service.CreateServiceAccountRequest{},
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      "The service account name is invalid",
		},
		{
			name: "empty name handled by MinIO",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetAddServiceAccountError(400, "The service account name is invalid")
			},
			requestBody: service.CreateServiceAccountRequest{
				Name: "",
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      "The service account name is invalid",
		},
		{
			name:      "whitespace-only name validation",
			setupMock: func(mock *minio.MockMinIOServer) {},
			requestBody: service.CreateServiceAccountRequest{
				Name: "   ",
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      "failed to create service account:",
		},
		{
			name:               "invalid JSON body",
			setupMock:          func(mock *minio.MockMinIOServer) {},
			requestBody:        "invalid json",
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      "Invalid request body",
		},
		{
			name: "MinIO server error",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetAddServiceAccountError(500, "Internal Server Error")
			},
			requestBody: service.CreateServiceAccountRequest{
				Name: "error-service-account",
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      "Internal Server Error",
		},
		{
			name: "MinIO authentication error",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetAddServiceAccountError(401, "Unauthorized")
			},
			requestBody: service.CreateServiceAccountRequest{
				Name: "auth-error-service-account",
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      "Unauthorized",
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
			addServiceAccountService := service.NewAddServiceAccountService(minioClient)

			// Create HTTP service
			testService := createTestService(t, nil, nil, addServiceAccountService)

			// Prepare request body
			var requestBody []byte
			if str, ok := tt.requestBody.(string); ok {
				requestBody = []byte(str)
			} else {
				requestBody, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/api/access-keys", bytes.NewReader(requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Add context with logger
			ctx := zerolog.New(zerolog.NewTestWriter(t)).WithContext(req.Context())
			req = req.WithContext(ctx)

			// Create response recorder
			w := httptest.NewRecorder()

			// Execute request
			testService.PostAccessKeysHandler(w, req)

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

				var response service.CreateServiceAccountResponse
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

func TestService_PostAccessKeysHandler_ContentTypeValidation(t *testing.T) {
	tests := []struct {
		name               string
		contentType        string
		requestBody        string
		expectedStatusCode int
		expectedError      string
	}{
		{
			name:               "missing content type",
			contentType:        "",
			requestBody:        `{"name": "test-service-account"}`,
			expectedStatusCode: http.StatusCreated, // Still works
		},
		{
			name:               "wrong content type",
			contentType:        "text/plain",
			requestBody:        `{"name": "test-service-account"}`,
			expectedStatusCode: http.StatusCreated, // Still works (Go is flexible)
		},
		{
			name:               "correct content type",
			contentType:        "application/json",
			requestBody:        `{"name": "test-service-account"}`,
			expectedStatusCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock MinIO server
			mockServer := minio.NewMockMinIOServer()
			defer mockServer.Close()

			scenarios := minio.TestScenarios{}
			mockServer.SetAddServiceAccountResponse(scenarios.SuccessfulAddServiceAccount())

			// Create MinIO client
			minioClient, err := mockServer.CreateMinIOClient()
			if err != nil {
				t.Fatalf("Failed to create MinIO client: %v", err)
			}

			// Create service dependencies
			addServiceAccountService := service.NewAddServiceAccountService(minioClient)

			// Create HTTP service
			testService := createTestService(t, nil, nil, addServiceAccountService)

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/api/access-keys", bytes.NewReader([]byte(tt.requestBody)))
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}

			// Add context with logger
			ctx := zerolog.New(zerolog.NewTestWriter(t)).WithContext(req.Context())
			req = req.WithContext(ctx)

			// Create response recorder
			w := httptest.NewRecorder()

			// Execute request
			testService.PostAccessKeysHandler(w, req)

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

func TestService_PostAccessKeysHandler_EdgeCases(t *testing.T) {
	tests := []struct {
		name               string
		setupMock          func(*minio.MockMinIOServer)
		setupRequest       func() *http.Request
		expectedStatusCode int
		expectedError      string
		validateResponse   func(t *testing.T, response *service.CreateServiceAccountResponse)
	}{
		{
			name:      "empty request body",
			setupMock: func(mock *minio.MockMinIOServer) {},
			setupRequest: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/api/access-keys", bytes.NewReader([]byte("")))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      "Invalid request body",
		},
		{
			name: "null JSON body",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetAddServiceAccountError(400, "The service account name is invalid")
			},
			setupRequest: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/api/access-keys", bytes.NewReader([]byte("null")))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      "The service account name is invalid",
		},
		{
			name: "empty JSON object handled by MinIO",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetAddServiceAccountError(400, "The service account name is invalid")
			},
			setupRequest: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/api/access-keys", bytes.NewReader([]byte("{}")))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      "The service account name is invalid",
		},
		{
			name: "request with extra fields",
			setupMock: func(mock *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				mock.SetAddServiceAccountResponse(scenarios.SuccessfulAddServiceAccount())
			},
			setupRequest: func() *http.Request {
				// JSON with extra fields that should be ignored
				body := `{
					"name": "extra-fields-test",
					"description": "Test with extra fields",
					"extraField": "should be ignored",
					"anotherExtra": 123,
					"nestedExtra": {
						"ignored": true
					}
				}`
				return httptest.NewRequest(http.MethodPost, "/api/access-keys", bytes.NewReader([]byte(body)))
			},
			expectedStatusCode: http.StatusCreated,
			validateResponse: func(t *testing.T, response *service.CreateServiceAccountResponse) {
				if response.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", response.AccessKey)
				}
				if response.Name != "extra-fields-test" {
					t.Errorf("Expected Name %q, got %q", "extra-fields-test", response.Name)
				}
				if response.Description != "Test with extra fields" {
					t.Errorf("Expected Description %q, got %q", "Test with extra fields", response.Description)
				}
			},
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
			addServiceAccountService := service.NewAddServiceAccountService(minioClient)

			// Create HTTP service
			testService := createTestService(t, nil, nil, addServiceAccountService)

			// Setup request
			req := tt.setupRequest()
			req.Header.Set("Content-Type", "application/json")

			// Add context with logger
			ctx := zerolog.New(zerolog.NewTestWriter(t)).WithContext(req.Context())
			req = req.WithContext(ctx)

			// Create response recorder
			w := httptest.NewRecorder()

			// Execute request
			testService.PostAccessKeysHandler(w, req)

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

				var response service.CreateServiceAccountResponse
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

// createTestService creates a Service instance for testing with mock dependencies
func createTestService(
	t *testing.T,
	getServerInfoService *service.GetServerInfoService,
	listAccessKeysService *service.ListAccessKeysService,
	addServiceAccountService *service.AddServiceAccountService,
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
		config:                   cfg,
		logger:                   logger,
		getServerInfoService:     getServerInfoService,
		listAccessKeysService:    listAccessKeysService,
		addServiceAccountService: addServiceAccountService,
	}
}

// TestService_PostAccessKeysHandler_Integration tests the handler through the full router
func TestService_PostAccessKeysHandler_Integration(t *testing.T) {
	// Setup mock MinIO server
	mockServer := minio.NewMockMinIOServer()
	defer mockServer.Close()

	scenarios := minio.TestScenarios{}
	mockServer.SetAddServiceAccountResponse(scenarios.SuccessfulAddServiceAccount())

	// Create MinIO client
	minioClient, err := mockServer.CreateMinIOClient()
	if err != nil {
		t.Fatalf("Failed to create MinIO client: %v", err)
	}

	// Create service dependencies
	addServiceAccountService := service.NewAddServiceAccountService(minioClient)

	// Create full service
	testService := createTestService(t, nil, nil, addServiceAccountService)

	// Create router and add our handler
	router := chi.NewRouter()
	router.Post("/api/access-keys", testService.PostAccessKeysHandler)

	// Create test server
	server := httptest.NewServer(router)
	defer server.Close()

	// Prepare request
	requestBody := service.CreateServiceAccountRequest{
		Name:        "integration-test",
		Description: "Integration test service account",
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Make request
	resp, err := http.Post(server.URL+"/api/access-keys", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Validate response
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type %q, got %q", "application/json", contentType)
	}

	var response service.CreateServiceAccountResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
		t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", response.AccessKey)
	}
	if response.SecretKey != "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" {
		t.Errorf("Expected SecretKey %q, got %q", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", response.SecretKey)
	}
	if response.Name != "integration-test" {
		t.Errorf("Expected Name %q, got %q", "integration-test", response.Name)
	}
	if response.Description != "Integration test service account" {
		t.Errorf("Expected Description %q, got %q", "Integration test service account", response.Description)
	}
}
