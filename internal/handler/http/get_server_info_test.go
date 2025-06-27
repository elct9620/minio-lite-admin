package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elct9620/minio-lite-admin/internal/testability/minio"
)

func TestService_GetServerInfoHandler(t *testing.T) {
	// Test cases using table-driven tests
	tests := []struct {
		name           string
		method         string
		setupMock      func(*minio.MockMinIOServer)
		expectedStatus int
		expectedBody   *ServerInfoResponse
		checkHeaders   map[string]string
	}{
		{
			name:   "successful server info request",
			method: http.MethodGet,
			setupMock: func(mockMinIO *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				mockMinIO.SetServerInfoResponse(scenarios.SuccessfulServerInfo())
			},
			expectedStatus: http.StatusOK,
			expectedBody: &ServerInfoResponse{
				Mode:         "standalone",
				Region:       "us-east-1",
				DeploymentID: "9b4b8c6f-1234-5678-9abc-123456789def",
			},
			checkHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name:   "distributed mode server info",
			method: http.MethodGet,
			setupMock: func(mockMinIO *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				mockMinIO.SetServerInfoResponse(scenarios.DistributedServerInfo())
			},
			expectedStatus: http.StatusOK,
			expectedBody: &ServerInfoResponse{
				Mode:         "distributed",
				Region:       "us-west-2",
				DeploymentID: "distributed-cluster-uuid-5678",
			},
			checkHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name:   "server info with empty region",
			method: http.MethodGet,
			setupMock: func(mockMinIO *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				mockMinIO.SetServerInfoResponse(scenarios.EmptyRegionServerInfo())
			},
			expectedStatus: http.StatusOK,
			expectedBody: &ServerInfoResponse{
				Mode:         "standalone",
				Region:       "",
				DeploymentID: "no-region-deployment-id",
			},
			checkHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name:   "MinIO server returns error",
			method: http.MethodGet,
			setupMock: func(mockMinIO *minio.MockMinIOServer) {
				// Use HTTP 400 (Bad Request) instead of 500 - it's non-retryable
				mockMinIO.SetServerInfoNonRetryableError(http.StatusBadRequest, "MinIO server error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil, // Expecting error response, not JSON
			checkHeaders: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
		},
		{
			name:   "POST method should also work",
			method: http.MethodPost,
			setupMock: func(mockMinIO *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				mockMinIO.SetServerInfoResponse(scenarios.SuccessfulServerInfo())
			},
			expectedStatus: http.StatusOK,
			expectedBody: &ServerInfoResponse{
				Mode:         "standalone",
				Region:       "us-east-1",
				DeploymentID: "9b4b8c6f-1234-5678-9abc-123456789def",
			},
			checkHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	// Run table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test service with mock MinIO
			svc, mockMinIO := testServiceWithMockMinIO()
			defer mockMinIO.Close()

			// Setup mock behavior
			if tt.setupMock != nil {
				tt.setupMock(mockMinIO)
			}

			// Create request
			req := httptest.NewRequest(tt.method, "/api/server-info", nil)
			
			// Create response recorder
			w := httptest.NewRecorder()

			// Call the handler
			svc.GetServerInfoHandler(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("GetServerInfoHandler() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			// Check headers
			for header, expectedValue := range tt.checkHeaders {
				if got := w.Header().Get(header); got != expectedValue {
					t.Errorf("GetServerInfoHandler() header %s = %v, want %v", header, got, expectedValue)
				}
			}

			// Check response body for successful cases
			if tt.expectedBody != nil {
				var response ServerInfoResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("GetServerInfoHandler() failed to unmarshal response: %v", err)
					return
				}

				if response.Mode != tt.expectedBody.Mode {
					t.Errorf("GetServerInfoHandler() mode = %v, want %v", response.Mode, tt.expectedBody.Mode)
				}

				if response.Region != tt.expectedBody.Region {
					t.Errorf("GetServerInfoHandler() region = %v, want %v", response.Region, tt.expectedBody.Region)
				}

				if response.DeploymentID != tt.expectedBody.DeploymentID {
					t.Errorf("GetServerInfoHandler() deploymentId = %v, want %v", response.DeploymentID, tt.expectedBody.DeploymentID)
				}
			}
		})
	}
}

func TestService_GetServerInfoHandler_ResponseFormat(t *testing.T) {
	// Test to ensure response is valid JSON and has correct structure
	svc, mockMinIO := testServiceWithMockMinIO()
	defer mockMinIO.Close()

	// Setup successful response
	scenarios := minio.TestScenarios{}
	mockMinIO.SetServerInfoResponse(scenarios.SuccessfulServerInfo())
	
	req := httptest.NewRequest(http.MethodGet, "/api/server-info", nil)
	w := httptest.NewRecorder()

	svc.GetServerInfoHandler(w, req)

	// Should return 200 OK
	if w.Code != http.StatusOK {
		t.Errorf("GetServerInfoHandler() status = %v, want %v", w.Code, http.StatusOK)
		return
	}

	// Check that response is valid JSON
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("GetServerInfoHandler() response is not valid JSON: %v", err)
		return
	}

	// Check that required fields exist
	expectedFields := []string{"mode", "region", "deploymentId"}
	for _, field := range expectedFields {
		if _, exists := response[field]; !exists {
			t.Errorf("GetServerInfoHandler() response missing field: %s", field)
		}
	}

	// Check field types
	if mode, ok := response["mode"].(string); !ok || mode == "" {
		t.Errorf("GetServerInfoHandler() mode field should be non-empty string, got: %v", response["mode"])
	}

	if _, ok := response["region"].(string); !ok {
		t.Errorf("GetServerInfoHandler() region field should be string, got: %v", response["region"])
	}

	if deploymentId, ok := response["deploymentId"].(string); !ok || deploymentId == "" {
		t.Errorf("GetServerInfoHandler() deploymentId field should be non-empty string, got: %v", response["deploymentId"])
	}
}

func TestService_GetServerInfoHandler_ErrorHandling(t *testing.T) {
	// Test error scenarios
	svc, mockMinIO := testServiceWithMockMinIO()
	defer mockMinIO.Close()

	// Setup non-retryable error response (HTTP 400 is not retried by madmin)
	mockMinIO.SetServerInfoNonRetryableError(http.StatusBadRequest, "Bad Request")
	
	req := httptest.NewRequest(http.MethodGet, "/api/server-info", nil)
	w := httptest.NewRecorder()

	svc.GetServerInfoHandler(w, req)

	// Should return 500 Internal Server Error (our handler converts MinIO errors to 500)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("GetServerInfoHandler() status = %v, want %v", w.Code, http.StatusInternalServerError)
	}

	// Should have error message in body
	body := w.Body.String()
	if body == "" {
		t.Error("GetServerInfoHandler() error response should have non-empty body")
	}
}