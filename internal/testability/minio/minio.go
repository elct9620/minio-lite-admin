package minio

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/elct9620/minio-lite-admin/internal/infra"
	"github.com/go-chi/chi/v5"
	"github.com/minio/madmin-go/v4"
)

// MockMinIOServer provides a mock MinIO server for testing
type MockMinIOServer struct {
	server          *httptest.Server
	responses       map[string]any
	serviceAccounts map[string]*ServiceAccountInfo // In-memory store for service accounts
}

// ServiceAccountInfo represents stored service account information
type ServiceAccountInfo struct {
	AccessKey   string          `json:"accessKey"`
	SecretKey   string          `json:"secretKey"`
	Name        string          `json:"name,omitempty"`
	Description string          `json:"description,omitempty"`
	Status      string          `json:"status"`
	Policy      json.RawMessage `json:"policy,omitempty"`
	Expiration  *time.Time      `json:"expiration,omitempty"`
	ParentUser  string          `json:"parentUser"`
	CreatedAt   time.Time       `json:"createdAt"`
}

// NewMockMinIOServer creates a new mock MinIO server
func NewMockMinIOServer() *MockMinIOServer {
	mock := &MockMinIOServer{
		responses:       make(map[string]any),
		serviceAccounts: make(map[string]*ServiceAccountInfo),
	}

	// Set default server info response
	mock.SetServerInfoResponse(ServerInfoResponse{
		Mode:         "standalone",
		Region:       "us-east-1",
		DeploymentID: "9b4b8c6f-1234-5678-9abc-123456789def",
	})

	// Create HTTP test server with chi router
	r := chi.NewRouter()

	// MinIO Admin API endpoints
	r.Route("/minio/admin", func(r chi.Router) {
		// Server info endpoints
		r.Get("/v4/info", mock.handleServerInfo)
		r.Get("/v3/info", mock.handleServerInfo)
		r.Get("/info", mock.handleServerInfo)

		// Access keys endpoints
		r.Get("/v4/list-users", mock.handleListUsers)
		r.Get("/v4/list-access-keys-bulk", mock.handleListAccessKeysBulk)

		// Service account endpoints
		r.Get("/v4/info-service-account", mock.handleInfoServiceAccount)
		r.Put("/v4/add-service-account", mock.handleAddServiceAccount)
		r.Post("/v4/add-service-account", mock.handleAddServiceAccount)
		r.Put("/v4/update-service-account", mock.handleUpdateServiceAccount)
		r.Post("/v4/update-service-account", mock.handleUpdateServiceAccount)
		r.Delete("/v4/delete-service-account", mock.handleDeleteServiceAccount)
	})

	// Add a catch-all handler for unhandled requests
	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("Not Found"))
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

// AddServiceAccountToStore adds a service account to the in-memory store
func (m *MockMinIOServer) AddServiceAccountToStore(accessKey, secretKey, name, description, status, parentUser string, policy json.RawMessage, expiration *time.Time) {
	m.serviceAccounts[accessKey] = &ServiceAccountInfo{
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		Name:        name,
		Description: description,
		Status:      status,
		Policy:      policy,
		Expiration:  expiration,
		ParentUser:  parentUser,
		CreatedAt:   time.Now(),
	}
}

// GetServiceAccountFromStore retrieves a service account from the in-memory store
func (m *MockMinIOServer) GetServiceAccountFromStore(accessKey string) (*ServiceAccountInfo, bool) {
	sa, exists := m.serviceAccounts[accessKey]
	return sa, exists
}

// UpdateServiceAccountInStore updates a service account in the in-memory store
func (m *MockMinIOServer) UpdateServiceAccountInStore(accessKey string, updates map[string]any) bool {
	sa, exists := m.serviceAccounts[accessKey]
	if !exists {
		return false
	}

	// Apply updates
	if name, ok := updates["name"]; ok {
		sa.Name = name.(string)
	}
	if description, ok := updates["description"]; ok {
		sa.Description = description.(string)
	}
	if status, ok := updates["status"]; ok {
		sa.Status = status.(string)
	}
	if policy, ok := updates["policy"]; ok {
		sa.Policy = policy.(json.RawMessage)
	}
	if secretKey, ok := updates["secretKey"]; ok {
		sa.SecretKey = secretKey.(string)
	}
	if expiration, ok := updates["expiration"]; ok {
		if exp, ok := expiration.(*time.Time); ok {
			sa.Expiration = exp
		}
	}

	return true
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
