package minio

import (
	"encoding/json"
	"net/http"

	"github.com/minio/madmin-go/v4"
)

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
