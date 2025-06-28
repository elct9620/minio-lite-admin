package minio

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/minio/madmin-go/v4"
)

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

// UpdateServiceAccountRequest represents the MinIO admin update service account request
type UpdateServiceAccountRequest struct {
	NewPolicy      json.RawMessage `json:"newPolicy,omitempty"`
	NewSecretKey   string          `json:"newSecretKey,omitempty"`
	NewStatus      string          `json:"newStatus,omitempty"`
	NewName        string          `json:"newName,omitempty"`
	NewDescription string          `json:"newDescription,omitempty"`
	NewExpiration  *time.Time      `json:"newExpiration,omitempty"`
}

// InfoServiceAccountResponse represents the response from InfoServiceAccount
type InfoServiceAccountResponse struct {
	ParentUser    string          `json:"parentUser"`
	AccountStatus string          `json:"accountStatus"`
	ImpliedPolicy bool            `json:"impliedPolicy"`
	Policy        json.RawMessage `json:"policy,omitempty"`
	Name          string          `json:"name,omitempty"`
	Description   string          `json:"description,omitempty"`
	Expiration    *time.Time      `json:"expiration,omitempty"`
}

// Service Account Creation Methods

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

// Service Account Update Methods

// SetUpdateServiceAccountError sets an error response for update service account requests
func (m *MockMinIOServer) SetUpdateServiceAccountError(statusCode int, message string) {
	m.responses["update-service-account-error"] = struct {
		StatusCode int
		Message    string
	}{
		StatusCode: statusCode,
		Message:    message,
	}
}

// SetUpdateServiceAccountSuccess sets a success response for update service account requests
func (m *MockMinIOServer) SetUpdateServiceAccountSuccess() {
	m.responses["update-service-account"] = true
}

// handleUpdateServiceAccount handles the MinIO admin update service account endpoint
func (m *MockMinIOServer) handleUpdateServiceAccount(w http.ResponseWriter, r *http.Request) {
	// Check if we should return an error
	if errorResponse, exists := m.responses["update-service-account-error"]; exists {
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

	var req UpdateServiceAccountRequest
	if err := json.Unmarshal(decryptedBody, &req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update service account in store
	updates := make(map[string]interface{})
	if req.NewPolicy != nil && len(req.NewPolicy) > 0 {
		updates["policy"] = req.NewPolicy
	}
	if req.NewSecretKey != "" {
		updates["secretKey"] = req.NewSecretKey
	}
	if req.NewStatus != "" {
		updates["status"] = req.NewStatus
	}
	if req.NewName != "" {
		updates["name"] = req.NewName
	}
	if req.NewDescription != "" {
		updates["description"] = req.NewDescription
	}
	if req.NewExpiration != nil {
		updates["expiration"] = req.NewExpiration
	}

	if !m.UpdateServiceAccountInStore(accessKey, updates) {
		http.Error(w, "Service account not found", http.StatusNotFound)
		return
	}

	// UpdateServiceAccount in madmin returns no response body
	w.WriteHeader(http.StatusNoContent)
}

// Service Account Deletion Methods

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

// Service Account Info Methods

// SetInfoServiceAccountResponse sets the response for info service account requests
func (m *MockMinIOServer) SetInfoServiceAccountResponse(accessKey string, response InfoServiceAccountResponse) {
	if m.responses["info-service-account"] == nil {
		m.responses["info-service-account"] = make(map[string]InfoServiceAccountResponse)
	}
	m.responses["info-service-account"].(map[string]InfoServiceAccountResponse)[accessKey] = response
}

// SetInfoServiceAccountError sets an error response for info service account requests
func (m *MockMinIOServer) SetInfoServiceAccountError(statusCode int, message string) {
	m.responses["info-service-account-error"] = struct {
		StatusCode int
		Message    string
	}{
		StatusCode: statusCode,
		Message:    message,
	}
}

// handleInfoServiceAccount handles the MinIO admin info service account endpoint
func (m *MockMinIOServer) handleInfoServiceAccount(w http.ResponseWriter, r *http.Request) {
	// Check if we should return an error
	if errorResponse, exists := m.responses["info-service-account-error"]; exists {
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

	var responseData InfoServiceAccountResponse

	// Check if we have a response for this specific access key
	if responses, exists := m.responses["info-service-account"]; exists {
		if infoResponses, ok := responses.(map[string]InfoServiceAccountResponse); ok {
			if info, found := infoResponses[accessKey]; found {
				responseData = info
			} else {
				// Try to get from service account store
				if saInfo, found := m.GetServiceAccountFromStore(accessKey); found {
					responseData = InfoServiceAccountResponse{
						ParentUser:    saInfo.ParentUser,
						AccountStatus: saInfo.Status,
						ImpliedPolicy: true,
						Policy:        saInfo.Policy,
						Name:          saInfo.Name,
						Description:   saInfo.Description,
						Expiration:    saInfo.Expiration,
					}
				} else {
					http.Error(w, "Service account not found", http.StatusNotFound)
					return
				}
			}
		}
	} else {
		// Try to get from service account store
		if saInfo, found := m.GetServiceAccountFromStore(accessKey); found {
			responseData = InfoServiceAccountResponse{
				ParentUser:    saInfo.ParentUser,
				AccountStatus: saInfo.Status,
				ImpliedPolicy: true,
				Policy:        saInfo.Policy,
				Name:          saInfo.Name,
				Description:   saInfo.Description,
				Expiration:    saInfo.Expiration,
			}
		} else {
			http.Error(w, "Service account not found", http.StatusNotFound)
			return
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
