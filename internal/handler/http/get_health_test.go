package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestService_GetHealthHandler(t *testing.T) {
	// Test cases using table-driven tests
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedBody   HealthResponse
		checkHeaders   map[string]string
	}{
		{
			name:           "successful health check with GET",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedBody: HealthResponse{
				Status:  "ok",
				Service: "minio-lite-admin",
			},
			checkHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name:           "health check with POST method (should still work)",
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
			expectedBody: HealthResponse{
				Status:  "ok",
				Service: "minio-lite-admin",
			},
			checkHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name:           "health check with PUT method (should still work)",
			method:         http.MethodPut,
			expectedStatus: http.StatusOK,
			expectedBody: HealthResponse{
				Status:  "ok",
				Service: "minio-lite-admin",
			},
			checkHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	// Create test service instance
	svc := testService()

	// Run table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest(tt.method, "/api/health", nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Call the handler
			svc.GetHealthHandler(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("GetHealthHandler() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			// Check headers
			for header, expectedValue := range tt.checkHeaders {
				if got := w.Header().Get(header); got != expectedValue {
					t.Errorf("GetHealthHandler() header %s = %v, want %v", header, got, expectedValue)
				}
			}

			// Check response body
			var response HealthResponse
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Errorf("GetHealthHandler() failed to unmarshal response: %v", err)
				return
			}

			if response.Status != tt.expectedBody.Status {
				t.Errorf("GetHealthHandler() status = %v, want %v", response.Status, tt.expectedBody.Status)
			}

			if response.Service != tt.expectedBody.Service {
				t.Errorf("GetHealthHandler() service = %v, want %v", response.Service, tt.expectedBody.Service)
			}
		})
	}
}

func TestService_GetHealthHandler_ResponseFormat(t *testing.T) {
	// Test to ensure response is valid JSON and has correct structure
	svc := testService()

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()

	svc.GetHealthHandler(w, req)

	// Check that response is valid JSON
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("GetHealthHandler() response is not valid JSON: %v", err)
		return
	}

	// Check that required fields exist
	expectedFields := []string{"status", "service"}
	for _, field := range expectedFields {
		if _, exists := response[field]; !exists {
			t.Errorf("GetHealthHandler() response missing field: %s", field)
		}
	}

	// Check field types
	if status, ok := response["status"].(string); !ok || status == "" {
		t.Errorf("GetHealthHandler() status field should be non-empty string, got: %v", response["status"])
	}

	if service, ok := response["service"].(string); !ok || service == "" {
		t.Errorf("GetHealthHandler() service field should be non-empty string, got: %v", response["service"])
	}
}

func TestService_GetHealthHandler_ContentType(t *testing.T) {
	// Dedicated test for content-type header
	svc := testService()

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()

	svc.GetHealthHandler(w, req)

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("GetHealthHandler() Content-Type = %v, want application/json", contentType)
	}
}
