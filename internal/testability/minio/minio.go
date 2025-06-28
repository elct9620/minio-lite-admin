package minio

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/elct9620/minio-lite-admin/internal/infra"
	"github.com/go-chi/chi/v5"
	"github.com/minio/madmin-go/v4"
)

// MockMinIOServer provides a mock MinIO server for testing
type MockMinIOServer struct {
	server    *httptest.Server
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
		r.Get("/v4/list-users", mock.handleListUsers)
		r.Get("/v4/list-access-keys-bulk", mock.handleListAccessKeysBulk)
		r.Put("/v4/add-service-account", mock.handleAddServiceAccount)
		r.Post("/v4/add-service-account", mock.handleAddServiceAccount) // Try POST as well
		r.Delete("/v4/delete-service-account", mock.handleDeleteServiceAccount)
	})

	// Add a catch-all handler for unhandled requests
	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
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

// Access Keys related structures and methods

// AccessKeysUsersResponse represents the users list response
type AccessKeysUsersResponse map[string]AccessKeysUserInfo

type AccessKeysUserInfo struct {
	Status string `json:"status"`
}

// AccessKeysBulkResponse represents the bulk access keys response
type AccessKeysBulkResponse map[string]AccessKeysListResponse

type AccessKeysListResponse struct {
	ServiceAccounts []AccessKeysServiceAccount `json:"serviceAccounts"`
	STSKeys         []AccessKeysServiceAccount `json:"stsKeys"`
}

type AccessKeysServiceAccount struct {
	ParentUser    string `json:"parentUser"`
	AccountStatus string `json:"accountStatus"`
	AccessKey     string `json:"accessKey"`
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	ImpliedPolicy bool   `json:"impliedPolicy"`
}

// SetUsersResponse sets the response for list users requests
func (m *MockMinIOServer) SetUsersResponse(response AccessKeysUsersResponse) {
	m.responses["users"] = response
}

// SetAccessKeysBulkResponse sets the response for bulk access keys requests
func (m *MockMinIOServer) SetAccessKeysBulkResponse(response AccessKeysBulkResponse) {
	m.responses["access-keys-bulk"] = response
}

// SetUsersError sets an error response for list users requests
func (m *MockMinIOServer) SetUsersError(statusCode int, message string) {
	m.responses["users-error"] = struct {
		StatusCode int
		Message    string
	}{
		StatusCode: statusCode,
		Message:    message,
	}
}

// SetAccessKeysBulkError sets an error response for bulk access keys requests
func (m *MockMinIOServer) SetAccessKeysBulkError(statusCode int, message string) {
	m.responses["access-keys-bulk-error"] = struct {
		StatusCode int
		Message    string
	}{
		StatusCode: statusCode,
		Message:    message,
	}
}

// handleListUsers handles the MinIO admin list users endpoint
func (m *MockMinIOServer) handleListUsers(w http.ResponseWriter, r *http.Request) {
	// Check if we should return an error
	if errorResponse, exists := m.responses["users-error"]; exists {
		if err, ok := errorResponse.(struct {
			StatusCode int
			Message    string
		}); ok {
			http.Error(w, err.Message, err.StatusCode)
			return
		}
	}

	var responseData AccessKeysUsersResponse

	// Return success response
	if response, exists := m.responses["users"]; exists {
		if usersResp, ok := response.(AccessKeysUsersResponse); ok {
			responseData = usersResp
		}
	} else {
		// Default response with test users
		responseData = AccessKeysUsersResponse{
			"minioadmin": {Status: "enabled"},
			"testuser":   {Status: "enabled"},
		}
	}

	// Encode to JSON
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	// Encrypt the response using the same secret key as the client
	encryptedData, err := madmin.EncryptData("minioadmin", jsonData)
	if err != nil {
		http.Error(w, "Failed to encrypt response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err := w.Write(encryptedData); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// handleListAccessKeysBulk handles the MinIO admin bulk access keys endpoint
func (m *MockMinIOServer) handleListAccessKeysBulk(w http.ResponseWriter, r *http.Request) {
	// Check if we should return an error
	if errorResponse, exists := m.responses["access-keys-bulk-error"]; exists {
		if err, ok := errorResponse.(struct {
			StatusCode int
			Message    string
		}); ok {
			http.Error(w, err.Message, err.StatusCode)
			return
		}
	}

	var responseData AccessKeysBulkResponse

	// Return success response
	if response, exists := m.responses["access-keys-bulk"]; exists {
		if accessKeysResp, ok := response.(AccessKeysBulkResponse); ok {
			responseData = accessKeysResp
		}
	} else {
		// Default response with test access keys
		responseData = AccessKeysBulkResponse{
			"minioadmin": {
				ServiceAccounts: []AccessKeysServiceAccount{
					{
						ParentUser:    "minioadmin",
						AccountStatus: "enabled",
						AccessKey:     "test-service-account-1",
						Name:          "test-account",
						Description:   "Test service account",
						ImpliedPolicy: true,
					},
				},
				STSKeys: []AccessKeysServiceAccount{},
			},
			"testuser": {
				ServiceAccounts: []AccessKeysServiceAccount{},
				STSKeys:         []AccessKeysServiceAccount{},
			},
		}
	}

	// Encode to JSON
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	// Encrypt the response using the same secret key as the client
	encryptedData, err := madmin.EncryptData("minioadmin", jsonData)
	if err != nil {
		http.Error(w, "Failed to encrypt response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err := w.Write(encryptedData); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// Test scenarios for access keys

// SuccessfulAccessKeys returns a typical successful access keys response
func (TestScenarios) SuccessfulAccessKeys() (AccessKeysUsersResponse, AccessKeysBulkResponse) {
	users := AccessKeysUsersResponse{
		"minioadmin": {Status: "enabled"},
		"testuser":   {Status: "enabled"},
		"readonly":   {Status: "disabled"},
	}

	accessKeys := AccessKeysBulkResponse{
		"minioadmin": {
			ServiceAccounts: []AccessKeysServiceAccount{
				{
					ParentUser:    "minioadmin",
					AccountStatus: "enabled",
					AccessKey:     "AKIAIOSFODNN7EXAMPLE",
					Name:          "admin-service-account",
					Description:   "Administrative service account",
					ImpliedPolicy: true,
				},
				{
					ParentUser:    "minioadmin",
					AccountStatus: "enabled",
					AccessKey:     "AKIAI44QH8DHBEXAMPLE",
					Name:          "backup-service-account",
					Description:   "Backup service account",
					ImpliedPolicy: false,
				},
			},
			STSKeys: []AccessKeysServiceAccount{
				{
					ParentUser:    "minioadmin",
					AccountStatus: "enabled",
					AccessKey:     "ASIAIOSFODNN7EXAMPLE",
					Name:          "",
					Description:   "",
					ImpliedPolicy: true,
				},
			},
		},
		"testuser": {
			ServiceAccounts: []AccessKeysServiceAccount{
				{
					ParentUser:    "testuser",
					AccountStatus: "enabled",
					AccessKey:     "AKIATEST123456789012",
					Name:          "test-account",
					Description:   "Test user service account",
					ImpliedPolicy: false,
				},
			},
			STSKeys: []AccessKeysServiceAccount{},
		},
		"readonly": {
			ServiceAccounts: []AccessKeysServiceAccount{},
			STSKeys:         []AccessKeysServiceAccount{},
		},
	}

	return users, accessKeys
}

// Service Account creation related structures and methods

// AddServiceAccountRequest represents the request to create a service account in MinIO
// This should match madmin.AddServiceAccountReq
type AddServiceAccountRequest struct {
	Policy      json.RawMessage `json:"policy,omitempty"`
	TargetUser  string          `json:"targetUser,omitempty"`
	AccessKey   string          `json:"accessKey,omitempty"`
	SecretKey   string          `json:"secretKey,omitempty"`
	Name        string          `json:"name,omitempty"`
	Description string          `json:"description,omitempty"`
	Expiration  *time.Time      `json:"expiration,omitempty"`
}

// AddServiceAccountCredentials represents the credentials part of the response
type AddServiceAccountCredentials struct {
	AccessKey    string     `json:"accessKey"`
	SecretKey    string     `json:"secretKey"`
	SessionToken string     `json:"sessionToken,omitempty"`
	Expiration   *time.Time `json:"expiration,omitempty"`
}

// AddServiceAccountResponse represents the response from creating a service account
// This should match the madmin.AddServiceAccountResp structure that the client expects
type AddServiceAccountResponse struct {
	Credentials AddServiceAccountCredentials `json:"credentials"`
}

// SetAddServiceAccountResponse sets the response for add service account requests
func (m *MockMinIOServer) SetAddServiceAccountResponse(response AddServiceAccountResponse) {
	m.responses["add-service-account"] = response
}

// SetAddServiceAccountError sets an error response for add service account requests
func (m *MockMinIOServer) SetAddServiceAccountError(statusCode int, message string) {
	m.responses["add-service-account-error"] = struct {
		StatusCode int
		Message    string
	}{
		StatusCode: statusCode,
		Message:    message,
	}
}

// handleAddServiceAccount handles the MinIO admin add service account endpoint
func (m *MockMinIOServer) handleAddServiceAccount(w http.ResponseWriter, r *http.Request) {
	// Check if we should return an error
	if errorResponse, exists := m.responses["add-service-account-error"]; exists {
		if err, ok := errorResponse.(struct {
			StatusCode int
			Message    string
		}); ok {
			http.Error(w, err.Message, err.StatusCode)
			return
		}
	}

	// Read and decrypt request body
	encryptedBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Decrypt the request using the same secret key as the client
	decryptedBody, err := madmin.DecryptData("minioadmin", bytes.NewReader(encryptedBody))
	if err != nil {
		http.Error(w, "Failed to decrypt request body", http.StatusBadRequest)
		return
	}

	var req AddServiceAccountRequest
	if err := json.Unmarshal(decryptedBody, &req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var responseData AddServiceAccountResponse

	// Return success response
	if response, exists := m.responses["add-service-account"]; exists {
		if addResp, ok := response.(AddServiceAccountResponse); ok {
			responseData = addResp
		}
	} else {
		// Default response - generate credentials if not provided
		accessKey := req.AccessKey
		secretKey := req.SecretKey

		if accessKey == "" {
			accessKey = "AKIAIOSFODNN7EXAMPLE"
		}
		if secretKey == "" {
			secretKey = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
		}

		credentials := AddServiceAccountCredentials{
			AccessKey:    accessKey,
			SecretKey:    secretKey,
			SessionToken: "",
		}

		// Only set expiration if it's not nil and not zero
		if req.Expiration != nil && !req.Expiration.IsZero() {
			credentials.Expiration = req.Expiration
		}

		responseData = AddServiceAccountResponse{
			Credentials: credentials,
		}
	}

	// Encode to JSON
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	// Encrypt the response using the same secret key as the client
	encryptedData, err := madmin.EncryptData("minioadmin", jsonData)
	if err != nil {
		http.Error(w, "Failed to encrypt response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err := w.Write(encryptedData); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// SuccessfulAddServiceAccount returns a typical successful add service account response
func (TestScenarios) SuccessfulAddServiceAccount() AddServiceAccountResponse {
	return AddServiceAccountResponse{
		Credentials: AddServiceAccountCredentials{
			AccessKey:    "AKIAIOSFODNN7EXAMPLE",
			SecretKey:    "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			SessionToken: "",
			Expiration:   nil, // No expiration
		},
	}
}

// CustomKeysAddServiceAccount returns a response with custom access and secret keys
func (TestScenarios) CustomKeysAddServiceAccount(accessKey, secretKey string) AddServiceAccountResponse {
	return AddServiceAccountResponse{
		Credentials: AddServiceAccountCredentials{
			AccessKey:    accessKey,
			SecretKey:    secretKey,
			SessionToken: "",
			Expiration:   nil,
		},
	}
}

// ExpiringAddServiceAccount returns a response with expiration
func (TestScenarios) ExpiringAddServiceAccount(expiration time.Time) AddServiceAccountResponse {
	return AddServiceAccountResponse{
		Credentials: AddServiceAccountCredentials{
			AccessKey:    "AKIAIOSFODNN7EXAMPLE",
			SecretKey:    "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			SessionToken: "",
			Expiration:   &expiration,
		},
	}
}

// Delete service account related structures and methods

// SetDeleteServiceAccountError sets an error response for delete service account requests
func (m *MockMinIOServer) SetDeleteServiceAccountError(statusCode int, message string) {
	m.responses["delete-service-account-error"] = struct {
		StatusCode int
		Message    string
	}{
		StatusCode: statusCode,
		Message:    message,
	}
}

// SetDeleteServiceAccountSuccess sets a success response for delete service account requests
func (m *MockMinIOServer) SetDeleteServiceAccountSuccess() {
	m.responses["delete-service-account"] = true
}

// handleDeleteServiceAccount handles the MinIO admin delete service account endpoint
func (m *MockMinIOServer) handleDeleteServiceAccount(w http.ResponseWriter, r *http.Request) {
	// Check if we should return an error
	if errorResponse, exists := m.responses["delete-service-account-error"]; exists {
		if err, ok := errorResponse.(struct {
			StatusCode int
			Message    string
		}); ok {
			http.Error(w, err.Message, err.StatusCode)
			return
		}
	}

	// Parse the accessKey parameter from query string
	accessKey := r.URL.Query().Get("accessKey")
	if accessKey == "" {
		http.Error(w, "Missing accessKey parameter", http.StatusBadRequest)
		return
	}

	// Return success response (empty body for delete operations)
	w.WriteHeader(http.StatusNoContent)
}

// Test scenarios for delete service account

// SuccessfulDeleteServiceAccount returns a successful delete operation
func (TestScenarios) SuccessfulDeleteServiceAccount() bool {
	return true
}
