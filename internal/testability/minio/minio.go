package minio

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/minio/madmin-go/v4"
	"github.com/elct9620/minio-lite-admin/internal/infra"
)

// MockMinIOServer provides a mock MinIO server for testing
type MockMinIOServer struct {
	server   *httptest.Server
	responses map[string]interface{}
}

// ServerInfoResponse represents the MinIO admin API server info response format
type ServerInfoResponse struct {
	Mode         string `json:"mode"`
	Region       string `json:"region"`
	DeploymentID string `json:"deploymentId"`
}

// NewMockMinIOServer creates a new mock MinIO server
func NewMockMinIOServer() *MockMinIOServer {
	mock := &MockMinIOServer{
		responses: make(map[string]interface{}),
	}

	// Set default responses
	mock.SetServerInfoResponse(ServerInfoResponse{
		Mode:         "standalone",
		Region:       "us-east-1",
		DeploymentID: "9b4b8c6f-1234-5678-9abc-123456789def",
	})

	// Create HTTP test server with chi router
	r := chi.NewRouter()
	
	// MinIO Admin API endpoints
	r.Route("/minio/admin", func(r chi.Router) {
		r.Get("/v4/info", mock.handleServerInfo)
		r.Get("/v3/info", mock.handleServerInfo)
		r.Get("/info", mock.handleServerInfo)
	})
	
	mock.server = httptest.NewServer(r)
	
	return mock
}

// Close shuts down the mock server
func (m *MockMinIOServer) Close() {
	m.server.Close()
}

// URL returns the mock server URL
func (m *MockMinIOServer) URL() string {
	return m.server.URL
}

// SetServerInfoResponse sets the response for server info requests
func (m *MockMinIOServer) SetServerInfoResponse(response ServerInfoResponse) {
	m.responses["server-info"] = response
}

// SetServerInfoError sets an error response for server info requests
func (m *MockMinIOServer) SetServerInfoError(statusCode int, message string) {
	m.responses["server-info-error"] = struct {
		StatusCode int
		Message    string
	}{
		StatusCode: statusCode,
		Message:    message,
	}
}

// SetServerInfoNonRetryableError sets a non-retryable error response for server info requests
// Uses HTTP status codes that madmin won't retry: 400, 401, 403, 404, etc.
//
// madmin retries these status codes (avoid in tests): 408, 429, 502, 503
// madmin does NOT retry these status codes (use in tests): 400, 401, 403, 404, 405, 409, 422, etc.
func (m *MockMinIOServer) SetServerInfoNonRetryableError(statusCode int, message string) {
	m.responses["server-info-error"] = struct {
		StatusCode int
		Message    string
	}{
		StatusCode: statusCode,
		Message:    message,
	}
}

// CreateMinIOClient creates a madmin.AdminClient configured to use this mock server
func (m *MockMinIOServer) CreateMinIOClient() (*madmin.AdminClient, error) {
	config := infra.MinIOConfig{
		URL:      m.server.URL,
		RootUser: "minioadmin",
		Password: "minioadmin",
	}

	return infra.NewMinIOClient(config)
}

// CreateMinIOClientWithTimeout creates a madmin.AdminClient with custom timeout for testing
func (m *MockMinIOServer) CreateMinIOClientWithTimeout(timeout time.Duration) (*madmin.AdminClient, error) {
	config := infra.MinIOConfig{
		URL:      m.server.URL,
		RootUser: "minioadmin",
		Password: "minioadmin",
	}

	return infra.NewMinIOClientWithTimeout(config, timeout)
}

// handleServerInfo handles the MinIO admin server info endpoint
func (m *MockMinIOServer) handleServerInfo(w http.ResponseWriter, r *http.Request) {
	// Check if we should return an error
	if errorResponse, exists := m.responses["server-info-error"]; exists {
		if err, ok := errorResponse.(struct {
			StatusCode int
			Message    string
		}); ok {
			http.Error(w, err.Message, err.StatusCode)
			// Don't delete error response - keep it for consistent behavior across retries
			return
		}
	}

	// Return success response
	if response, exists := m.responses["server-info"]; exists {
		if serverInfo, ok := response.(ServerInfoResponse); ok {
			// MinIO admin API returns server info in a specific format
			// We need to match the expected response structure for madmin.InfoMessage
			minioResponse := map[string]interface{}{
				"mode":         serverInfo.Mode,
				"region":       serverInfo.Region,
				"deploymentId": serverInfo.DeploymentID,
				"platform":     "linux",
				"runtime":      "go1.21.0",
				"servers": []map[string]interface{}{
					{
						"endpoint": m.server.URL,
						"uptime":   3600,
						"version":  "minio-test",
						"commitID": "test-commit",
						"network":  map[string]interface{}{},
						"drives":   []interface{}{},
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(minioResponse); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			}
			return
		}
	}

	// Default fallback response
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}


// TestScenarios provides pre-configured test scenarios
type TestScenarios struct{}

// SuccessfulServerInfo returns a typical successful server info response
func (TestScenarios) SuccessfulServerInfo() ServerInfoResponse {
	return ServerInfoResponse{
		Mode:         "standalone",
		Region:       "us-east-1",
		DeploymentID: "9b4b8c6f-1234-5678-9abc-123456789def",
	}
}

// DistributedServerInfo returns a distributed mode server info response
func (TestScenarios) DistributedServerInfo() ServerInfoResponse {
	return ServerInfoResponse{
		Mode:         "distributed",
		Region:       "us-west-2",
		DeploymentID: "distributed-cluster-uuid-5678",
	}
}

// EmptyRegionServerInfo returns server info with empty region
func (TestScenarios) EmptyRegionServerInfo() ServerInfoResponse {
	return ServerInfoResponse{
		Mode:         "standalone",
		Region:       "",
		DeploymentID: "no-region-deployment-id",
	}
}