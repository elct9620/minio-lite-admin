package http

import (
	"bytes"
	"context"
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
	"github.com/minio/madmin-go/v4"
	"github.com/rs/zerolog"
)

func TestService_PutAccessKeysHandler(t *testing.T) {
	tests := []struct {
		name               string
		setupMock          func(*minio.MockMinIOServer)
		accessKey          string
		requestBody        interface{}
		expectedStatusCode int
		expectedError      string
		validateResponse   func(t *testing.T, response *service.UpdateServiceAccountResponse)
	}{
		{
			name: "successful update with policy",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetUpdateServiceAccountSuccess()
			},
			accessKey: "AKIAIOSFODNN7EXAMPLE",
			requestBody: map[string]interface{}{
				"newPolicy": `{"Version": "2012-10-17", "Statement": [{"Effect": "Allow", "Action": "s3:GetObject", "Resource": "arn:aws:s3:::bucket/*"}]}`,
			},
			expectedStatusCode: http.StatusOK,
			validateResponse: func(t *testing.T, response *service.UpdateServiceAccountResponse) {
				if response.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", response.AccessKey)
				}
				if response.Message != "Service account updated successfully" {
					t.Errorf("Expected Message %q, got %q", "Service account updated successfully", response.Message)
				}
			},
		},
		{
			name: "successful update with secret key rotation",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetUpdateServiceAccountSuccess()
			},
			accessKey: "AKIAIOSFODNN7EXAMPLE",
			requestBody: map[string]interface{}{
				"newSecretKey": "newSecretKey123456789012345",
			},
			expectedStatusCode: http.StatusOK,
			validateResponse: func(t *testing.T, response *service.UpdateServiceAccountResponse) {
				if response.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", response.AccessKey)
				}
				if response.Message != "Service account updated successfully" {
					t.Errorf("Expected Message %q, got %q", "Service account updated successfully", response.Message)
				}
			},
		},
		{
			name: "successful update with status change",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetUpdateServiceAccountSuccess()
			},
			accessKey: "AKIAIOSFODNN7EXAMPLE",
			requestBody: map[string]interface{}{
				"newStatus": string(madmin.AccountDisabled),
			},
			expectedStatusCode: http.StatusOK,
			validateResponse: func(t *testing.T, response *service.UpdateServiceAccountResponse) {
				if response.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", response.AccessKey)
				}
				if response.Message != "Service account updated successfully" {
					t.Errorf("Expected Message %q, got %q", "Service account updated successfully", response.Message)
				}
			},
		},
		{
			name: "successful update with name and description",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetUpdateServiceAccountSuccess()
			},
			accessKey: "AKIAIOSFODNN7EXAMPLE",
			requestBody: map[string]interface{}{
				"newName":        "Updated Service Account",
				"newDescription": "Updated description",
			},
			expectedStatusCode: http.StatusOK,
			validateResponse: func(t *testing.T, response *service.UpdateServiceAccountResponse) {
				if response.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", response.AccessKey)
				}
				if response.Message != "Service account updated successfully" {
					t.Errorf("Expected Message %q, got %q", "Service account updated successfully", response.Message)
				}
			},
		},
		{
			name: "successful update with expiration",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetUpdateServiceAccountSuccess()
			},
			accessKey: "AKIAIOSFODNN7EXAMPLE",
			requestBody: map[string]interface{}{
				"newExpiration": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			},
			expectedStatusCode: http.StatusOK,
			validateResponse: func(t *testing.T, response *service.UpdateServiceAccountResponse) {
				if response.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", response.AccessKey)
				}
				if response.Message != "Service account updated successfully" {
					t.Errorf("Expected Message %q, got %q", "Service account updated successfully", response.Message)
				}
			},
		},
		{
			name: "successful update with all fields",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetUpdateServiceAccountSuccess()
			},
			accessKey: "AKIAIOSFODNN7EXAMPLE",
			requestBody: map[string]interface{}{
				"newPolicy":      `{"Version": "2012-10-17", "Statement": [{"Effect": "Allow", "Action": "s3:*", "Resource": "*"}]}`,
				"newSecretKey":   "newSecretKey123456789012345",
				"newStatus":      string(madmin.AccountEnabled),
				"newName":        "Full Update Service Account",
				"newDescription": "Fully updated description",
				"newExpiration":  time.Now().Add(48 * time.Hour).Format(time.RFC3339),
			},
			expectedStatusCode: http.StatusOK,
			validateResponse: func(t *testing.T, response *service.UpdateServiceAccountResponse) {
				if response.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", response.AccessKey)
				}
				if response.Message != "Service account updated successfully" {
					t.Errorf("Expected Message %q, got %q", "Service account updated successfully", response.Message)
				}
			},
		},
		{
			name:               "missing access key parameter",
			setupMock:          func(mock *minio.MockMinIOServer) {},
			accessKey:          "",
			requestBody:        map[string]interface{}{},
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      "Access key is required",
		},
		{
			name:               "empty access key parameter",
			setupMock:          func(mock *minio.MockMinIOServer) {},
			accessKey:          "   ",
			requestBody:        map[string]interface{}{},
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      "Access key is required",
		},
		{
			name:               "invalid JSON body",
			setupMock:          func(mock *minio.MockMinIOServer) {},
			accessKey:          "AKIAIOSFODNN7EXAMPLE",
			requestBody:        "invalid-json",
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      "Invalid request body",
		},
		{
			name: "MinIO server error",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetUpdateServiceAccountError(500, "Internal Server Error")
			},
			accessKey: "AKIAIOSFODNN7EXAMPLE",
			requestBody: map[string]interface{}{
				"newPolicy": `{"Version": "2012-10-17", "Statement": []}`,
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      "Failed to update access key",
		},
		{
			name: "MinIO authentication error",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetUpdateServiceAccountError(401, "Unauthorized")
			},
			accessKey: "AKIAIOSFODNN7EXAMPLE",
			requestBody: map[string]interface{}{
				"newPolicy": `{"Version": "2012-10-17", "Statement": []}`,
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      "Failed to update access key",
		},
		{
			name: "MinIO not found error",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetUpdateServiceAccountError(404, "Service account not found")
			},
			accessKey: "NONEXISTENT",
			requestBody: map[string]interface{}{
				"newPolicy": `{"Version": "2012-10-17", "Statement": []}`,
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      "Failed to update access key",
		},
		{
			name: "MinIO forbidden error",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetUpdateServiceAccountError(403, "Forbidden")
			},
			accessKey: "AKIAIOSFODNN7EXAMPLE",
			requestBody: map[string]interface{}{
				"newPolicy": `{"Version": "2012-10-17", "Statement": []}`,
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      "Failed to update access key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock MinIO server
			mockServer := minio.NewMockMinIOServer()
			defer mockServer.Close()

			// Pre-populate service account in store for update tests
			if tt.accessKey != "" && tt.accessKey != "NONEXISTENT" {
				mockServer.AddServiceAccountToStore(
					tt.accessKey,
					"existingSecretKey",
					"Original Name",
					"Original Description",
					"enabled",
					"admin",
					nil, // no policy
					nil, // no expiration
				)
			}

			tt.setupMock(mockServer)

			// Create MinIO client
			minioClient, err := mockServer.CreateMinIOClient()
			if err != nil {
				t.Fatalf("Failed to create MinIO client: %v", err)
			}

			// Create service dependencies
			updateServiceAccountService := service.NewUpdateServiceAccountService(minioClient)

			// Create HTTP service
			testService := createTestServiceForUpdate(t, updateServiceAccountService)

			// Prepare request body
			var body bytes.Buffer
			if tt.requestBody != nil {
				if str, ok := tt.requestBody.(string); ok {
					// Handle invalid JSON case
					body.WriteString(str)
				} else {
					if err := json.NewEncoder(&body).Encode(tt.requestBody); err != nil {
						t.Fatalf("Failed to encode request body: %v", err)
					}
				}
			}

			// Create request with Chi context for URL parameters
			req := httptest.NewRequest(http.MethodPut, "/api/access-keys/test", &body)

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
			testService.PutAccessKeysHandler(w, req)

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

				var response service.UpdateServiceAccountResponse
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

func TestService_PutAccessKeysHandler_UrlParameterEdgeCases(t *testing.T) {
	tests := []struct {
		name               string
		setupRequest       func() *http.Request
		expectedStatusCode int
		expectedError      string
	}{
		{
			name: "URL parameter with special characters",
			setupRequest: func() *http.Request {
				body := bytes.NewBufferString(`{"newName": "Test Update"}`)
				req := httptest.NewRequest(http.MethodPut, "/api/access-keys/test", body)
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
				body := bytes.NewBufferString(`{"newName": "Test Update"}`)
				req := httptest.NewRequest(http.MethodPut, "/api/access-keys/test", body)
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
				body := bytes.NewBufferString(`{"newName": "Test Update"}`)
				req := httptest.NewRequest(http.MethodPut, "/api/access-keys/test", body)
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

			// Pre-populate service accounts for edge case tests
			mockServer.AddServiceAccountToStore(
				"AKIA TEST",
				"existingSecretKey",
				"Original Name",
				"Original Description",
				"enabled",
				"admin",
				nil, // no policy
				nil, // no expiration
			)
			mockServer.AddServiceAccountToStore(
				"AKIATEST",
				"existingSecretKey",
				"Original Name",
				"Original Description",
				"enabled",
				"admin",
				nil, // no policy
				nil, // no expiration
			)

			mockServer.SetUpdateServiceAccountSuccess()

			// Create MinIO client
			minioClient, err := mockServer.CreateMinIOClient()
			if err != nil {
				t.Fatalf("Failed to create MinIO client: %v", err)
			}

			// Create service dependencies
			updateServiceAccountService := service.NewUpdateServiceAccountService(minioClient)

			// Create HTTP service
			testService := createTestServiceForUpdate(t, updateServiceAccountService)

			// Setup request
			req := tt.setupRequest()

			// Add context with logger
			ctx := zerolog.New(zerolog.NewTestWriter(t)).WithContext(req.Context())
			req = req.WithContext(ctx)

			// Create response recorder
			w := httptest.NewRecorder()

			// Execute request
			testService.PutAccessKeysHandler(w, req)

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

// createTestServiceForUpdate creates a Service instance for testing update functionality
func createTestServiceForUpdate(
	t *testing.T,
	updateServiceAccountService *service.UpdateServiceAccountService,
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
		updateServiceAccountService: updateServiceAccountService,
	}
}

// TestService_PutAccessKeysHandler_Integration tests the handler through the full router
func TestService_PutAccessKeysHandler_Integration(t *testing.T) {
	// Setup mock MinIO server
	mockServer := minio.NewMockMinIOServer()
	defer mockServer.Close()

	// Pre-populate service account for integration test
	mockServer.AddServiceAccountToStore(
		"AKIAIOSFODNN7EXAMPLE",
		"existingSecretKey",
		"Original Name",
		"Original Description",
		"enabled",
		"admin",
		nil, // no policy
		nil, // no expiration
	)

	mockServer.SetUpdateServiceAccountSuccess()

	// Create MinIO client
	minioClient, err := mockServer.CreateMinIOClient()
	if err != nil {
		t.Fatalf("Failed to create MinIO client: %v", err)
	}

	// Create service dependencies
	updateServiceAccountService := service.NewUpdateServiceAccountService(minioClient)

	// Create full service
	testService := createTestServiceForUpdate(t, updateServiceAccountService)

	// Create router and add our handler
	router := chi.NewRouter()
	router.Put("/api/access-keys/{accessKey}", testService.PutAccessKeysHandler)

	// Create test server
	server := httptest.NewServer(router)
	defer server.Close()

	// Prepare request body
	requestBody := map[string]interface{}{
		"newPolicy": `{"Version": "2012-10-17", "Statement": [{"Effect": "Allow", "Action": "s3:GetObject", "Resource": "*"}]}`,
		"newName":   "Integration Test Service Account",
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Make request
	req, err := http.NewRequest(http.MethodPut, server.URL+"/api/access-keys/AKIAIOSFODNN7EXAMPLE", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Validate response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type %q, got %q", "application/json", contentType)
	}

	var response service.UpdateServiceAccountResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
		t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", response.AccessKey)
	}
	if response.Message != "Service account updated successfully" {
		t.Errorf("Expected Message %q, got %q", "Service account updated successfully", response.Message)
	}
}