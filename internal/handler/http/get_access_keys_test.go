package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elct9620/minio-lite-admin/internal/service"
	"github.com/elct9620/minio-lite-admin/internal/testability/minio"
)

func TestService_GetAccessKeysHandler(t *testing.T) {
	// Test cases using table-driven tests
	tests := []struct {
		name           string
		method         string
		queryParams    string
		setupMock      func(*minio.MockMinIOServer)
		expectedStatus int
		expectedFields []string
		checkHeaders   map[string]string
	}{
		{
			name:        "successful access keys request - all types",
			method:      http.MethodGet,
			queryParams: "",
			setupMock: func(mockMinIO *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				users, accessKeys := scenarios.SuccessfulAccessKeys()
				mockMinIO.SetUsersResponse(users)
				mockMinIO.SetAccessKeysBulkResponse(accessKeys)
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"accessKeys", "total"},
			checkHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name:        "successful access keys request - users only",
			method:      http.MethodGet,
			queryParams: "?type=users",
			setupMock: func(mockMinIO *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				users, accessKeys := scenarios.SuccessfulAccessKeys()
				mockMinIO.SetUsersResponse(users)
				mockMinIO.SetAccessKeysBulkResponse(accessKeys)
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"accessKeys", "total"},
			checkHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name:        "successful access keys request - service accounts only",
			method:      http.MethodGet,
			queryParams: "?type=serviceAccounts",
			setupMock: func(mockMinIO *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				users, accessKeys := scenarios.SuccessfulAccessKeys()
				mockMinIO.SetUsersResponse(users)
				mockMinIO.SetAccessKeysBulkResponse(accessKeys)
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"accessKeys", "total"},
			checkHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name:        "successful access keys request - STS only",
			method:      http.MethodGet,
			queryParams: "?type=sts",
			setupMock: func(mockMinIO *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				users, accessKeys := scenarios.SuccessfulAccessKeys()
				mockMinIO.SetUsersResponse(users)
				mockMinIO.SetAccessKeysBulkResponse(accessKeys)
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"accessKeys", "total"},
			checkHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name:        "successful access keys request - filter by user",
			method:      http.MethodGet,
			queryParams: "?user=testuser",
			setupMock: func(mockMinIO *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				users, accessKeys := scenarios.SuccessfulAccessKeys()
				mockMinIO.SetUsersResponse(users)
				mockMinIO.SetAccessKeysBulkResponse(accessKeys)
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"accessKeys", "total"},
			checkHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name:           "invalid type parameter",
			method:         http.MethodGet,
			queryParams:    "?type=invalid",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedFields: nil,
			checkHeaders: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
		},
		{
			name:        "MinIO users API returns error",
			method:      http.MethodGet,
			queryParams: "",
			setupMock: func(mockMinIO *minio.MockMinIOServer) {
				mockMinIO.SetUsersError(http.StatusBadRequest, "Bad Request")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedFields: nil,
			checkHeaders: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
		},
		{
			name:        "MinIO access keys bulk API returns error",
			method:      http.MethodGet,
			queryParams: "",
			setupMock: func(mockMinIO *minio.MockMinIOServer) {
				// Set valid users response but error for access keys bulk
				scenarios := minio.TestScenarios{}
				users, _ := scenarios.SuccessfulAccessKeys()
				mockMinIO.SetUsersResponse(users)
				mockMinIO.SetAccessKeysBulkError(http.StatusBadRequest, "Bad Request")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedFields: nil,
			checkHeaders: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
		},
		{
			name:        "POST method should also work",
			method:      http.MethodPost,
			queryParams: "",
			setupMock: func(mockMinIO *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				users, accessKeys := scenarios.SuccessfulAccessKeys()
				mockMinIO.SetUsersResponse(users)
				mockMinIO.SetAccessKeysBulkResponse(accessKeys)
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"accessKeys", "total"},
			checkHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	// Run table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test service with mock MinIO
			svc, mockMinIO := testServiceWithMockMinIOAndAccessKeys()
			defer mockMinIO.Close()

			// Setup mock behavior
			if tt.setupMock != nil {
				tt.setupMock(mockMinIO)
			}

			// Create request
			req := httptest.NewRequest(tt.method, "/api/access-keys"+tt.queryParams, nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Call the handler
			svc.GetAccessKeysHandler(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("GetAccessKeysHandler() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			// Check headers
			for header, expectedValue := range tt.checkHeaders {
				if got := w.Header().Get(header); got != expectedValue {
					t.Errorf("GetAccessKeysHandler() header %s = %v, want %v", header, got, expectedValue)
				}
			}

			// Check response body for successful cases
			if tt.expectedFields != nil {
				var response map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("GetAccessKeysHandler() failed to unmarshal response: %v", err)
					return
				}

				// Check that all expected fields exist
				for _, field := range tt.expectedFields {
					if _, exists := response[field]; !exists {
						t.Errorf("GetAccessKeysHandler() response missing field: %s", field)
					}
				}

				// Check specific field types
				if accessKeysField, ok := response["accessKeys"]; ok {
					if _, isArray := accessKeysField.([]interface{}); !isArray {
						t.Errorf("GetAccessKeysHandler() accessKeys should be array, got: %T", accessKeysField)
					}
				}

				if totalField, ok := response["total"]; ok {
					if _, isFloat64 := totalField.(float64); !isFloat64 {
						t.Errorf("GetAccessKeysHandler() total should be number, got: %T", totalField)
					}
				}
			}
		})
	}
}

func TestService_GetAccessKeysHandler_ResponseFormat(t *testing.T) {
	// Test to ensure response has correct structure and values
	svc, mockMinIO := testServiceWithMockMinIOAndAccessKeys()
	defer mockMinIO.Close()

	// Setup successful response
	scenarios := minio.TestScenarios{}
	users, accessKeys := scenarios.SuccessfulAccessKeys()
	mockMinIO.SetUsersResponse(users)
	mockMinIO.SetAccessKeysBulkResponse(accessKeys)

	req := httptest.NewRequest(http.MethodGet, "/api/access-keys", nil)
	w := httptest.NewRecorder()

	svc.GetAccessKeysHandler(w, req)

	// Should return 200 OK
	if w.Code != http.StatusOK {
		t.Errorf("GetAccessKeysHandler() status = %v, want %v", w.Code, http.StatusOK)
		return
	}

	// Check that response is valid JSON
	var response service.ListAccessKeysResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("GetAccessKeysHandler() response is not valid JSON: %v", err)
		return
	}

	// Check numeric fields are reasonable (non-negative)
	if response.Total < 0 {
		t.Error("GetAccessKeysHandler() total should be non-negative")
	}

	// Check that accessKeys is present and properly structured
	if response.AccessKeys == nil {
		t.Error("GetAccessKeysHandler() accessKeys should be an array, not nil")
	}

	// Check that we have some access keys in the test scenario
	if response.Total == 0 {
		t.Error("GetAccessKeysHandler() should return some access keys in test scenario")
	}

	// Verify access key structure if any exist
	if len(response.AccessKeys) > 0 {
		accessKey := response.AccessKeys[0]
		if accessKey.AccessKey == "" {
			t.Error("GetAccessKeysHandler() access key should have non-empty AccessKey")
		}
		if accessKey.ParentUser == "" {
			t.Error("GetAccessKeysHandler() access key should have non-empty ParentUser")
		}
		if accessKey.Type == "" {
			t.Error("GetAccessKeysHandler() access key should have non-empty Type")
		}
	}
}

func TestService_GetAccessKeysHandler_ErrorHandling(t *testing.T) {
	// Test error scenarios
	svc, mockMinIO := testServiceWithMockMinIOAndAccessKeys()
	defer mockMinIO.Close()

	// Setup non-retryable error response
	mockMinIO.SetUsersError(http.StatusBadRequest, "Bad Request")

	req := httptest.NewRequest(http.MethodGet, "/api/access-keys", nil)
	w := httptest.NewRecorder()

	svc.GetAccessKeysHandler(w, req)

	// Should return 500 Internal Server Error
	if w.Code != http.StatusInternalServerError {
		t.Errorf("GetAccessKeysHandler() status = %v, want %v", w.Code, http.StatusInternalServerError)
	}

	// Should have error message in body
	body := w.Body.String()
	if body == "" {
		t.Error("GetAccessKeysHandler() error response should have non-empty body")
	}
}
