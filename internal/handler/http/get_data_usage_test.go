package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elct9620/minio-lite-admin/internal/testability/minio"
)

func TestService_GetDataUsageHandler(t *testing.T) {
	// Test cases using table-driven tests
	tests := []struct {
		name           string
		method         string
		setupMock      func(*minio.MockMinIOServer)
		expectedStatus int
		expectedFields []string
		checkHeaders   map[string]string
	}{
		{
			name:   "successful data usage request",
			method: http.MethodGet,
			setupMock: func(mockMinIO *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				mockMinIO.SetServerInfoResponse(scenarios.SuccessfulServerInfo())
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"totalCapacity", "totalFreeCapacity", "totalUsedCapacity", "usagePercentage", "onlineDisks", "offlineDisks", "healingDisks", "poolsCount", "objectsCount", "bucketsCount", "diskDetails"},
			checkHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name:   "distributed mode data usage",
			method: http.MethodGet,
			setupMock: func(mockMinIO *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				mockMinIO.SetServerInfoResponse(scenarios.DistributedServerInfo())
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"totalCapacity", "totalFreeCapacity", "totalUsedCapacity", "usagePercentage", "onlineDisks", "offlineDisks", "healingDisks", "poolsCount", "objectsCount", "bucketsCount", "diskDetails"},
			checkHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name:   "MinIO server returns error",
			method: http.MethodGet,
			setupMock: func(mockMinIO *minio.MockMinIOServer) {
				mockMinIO.SetServerInfoNonRetryableError(http.StatusBadRequest, "Bad Request")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedFields: nil,
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
			expectedFields: []string{"totalCapacity", "totalFreeCapacity", "totalUsedCapacity", "usagePercentage", "onlineDisks", "offlineDisks", "healingDisks", "poolsCount", "objectsCount", "bucketsCount", "diskDetails"},
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
			req := httptest.NewRequest(tt.method, "/api/data-usage", nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Call the handler
			svc.GetDataUsageHandler(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("GetDataUsageHandler() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			// Check headers
			for header, expectedValue := range tt.checkHeaders {
				if got := w.Header().Get(header); got != expectedValue {
					t.Errorf("GetDataUsageHandler() header %s = %v, want %v", header, got, expectedValue)
				}
			}

			// Check response body for successful cases
			if tt.expectedFields != nil {
				var response map[string]any
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("GetDataUsageHandler() failed to unmarshal response: %v", err)
					return
				}

				// Check that all expected fields exist
				for _, field := range tt.expectedFields {
					if _, exists := response[field]; !exists {
						t.Errorf("GetDataUsageHandler() response missing field: %s", field)
					}
				}

				// Check specific field types
				if capacityField, ok := response["totalCapacity"]; ok {
					if _, isFloat64 := capacityField.(float64); !isFloat64 {
						t.Errorf("GetDataUsageHandler() totalCapacity should be number, got: %T", capacityField)
					}
				}

				if usageField, ok := response["usagePercentage"]; ok {
					if _, isFloat64 := usageField.(float64); !isFloat64 {
						t.Errorf("GetDataUsageHandler() usagePercentage should be number, got: %T", usageField)
					}
				}

				if disksField, ok := response["onlineDisks"]; ok {
					if _, isFloat64 := disksField.(float64); !isFloat64 {
						t.Errorf("GetDataUsageHandler() onlineDisks should be number, got: %T", disksField)
					}
				}

				if detailsField, ok := response["diskDetails"]; ok {
					if _, isArray := detailsField.([]any); !isArray {
						t.Errorf("GetDataUsageHandler() diskDetails should be array, got: %T", detailsField)
					}
				}
			}
		})
	}
}

func TestService_GetDataUsageHandler_ResponseFormat(t *testing.T) {
	// Test to ensure response has correct structure and values
	svc, mockMinIO := testServiceWithMockMinIO()
	defer mockMinIO.Close()

	// Setup successful response
	scenarios := minio.TestScenarios{}
	mockMinIO.SetServerInfoResponse(scenarios.SuccessfulServerInfo())

	req := httptest.NewRequest(http.MethodGet, "/api/data-usage", nil)
	w := httptest.NewRecorder()

	svc.GetDataUsageHandler(w, req)

	// Should return 200 OK
	if w.Code != http.StatusOK {
		t.Errorf("GetDataUsageHandler() status = %v, want %v", w.Code, http.StatusOK)
		return
	}

	// Check that response is valid JSON
	var response DataUsageResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("GetDataUsageHandler() response is not valid JSON: %v", err)
		return
	}

	// Check numeric fields are reasonable (TotalCapacity is uint64, so always >= 0)
	// Just check it's a valid value by ensuring the field exists
	_ = response.TotalCapacity

	if response.UsagePercentage < 0 || response.UsagePercentage > 100 {
		t.Errorf("GetDataUsageHandler() usagePercentage should be 0-100, got: %f", response.UsagePercentage)
	}

	if response.OnlineDisks < 0 {
		t.Error("GetDataUsageHandler() onlineDisks should be non-negative")
	}

	// Check that diskDetails is present and properly structured
	if response.DiskDetails == nil {
		t.Error("GetDataUsageHandler() diskDetails should be an array, not nil")
	}
}

func TestService_GetDataUsageHandler_ErrorHandling(t *testing.T) {
	// Test error scenarios
	svc, mockMinIO := testServiceWithMockMinIO()
	defer mockMinIO.Close()

	// Setup non-retryable error response
	mockMinIO.SetServerInfoNonRetryableError(http.StatusInternalServerError, "Internal Server Error")

	req := httptest.NewRequest(http.MethodGet, "/api/data-usage", nil)
	w := httptest.NewRecorder()

	svc.GetDataUsageHandler(w, req)

	// Should return 500 Internal Server Error
	if w.Code != http.StatusInternalServerError {
		t.Errorf("GetDataUsageHandler() status = %v, want %v", w.Code, http.StatusInternalServerError)
	}

	// Should have error message in body
	body := w.Body.String()
	if body == "" {
		t.Error("GetDataUsageHandler() error response should have non-empty body")
	}
}
