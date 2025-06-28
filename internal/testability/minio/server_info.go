package minio

import (
	"encoding/json"
	"net/http"
)

// ServerInfoResponse represents the MinIO admin API server info response format
type ServerInfoResponse struct {
	Mode         string `json:"mode"`
	Region       string `json:"region"`
	DeploymentID string `json:"deploymentId"`
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
