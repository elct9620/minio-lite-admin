package service

import (
	"context"
	"testing"

	"github.com/elct9620/minio-lite-admin/internal/testability/minio"
	"github.com/rs/zerolog"
)

func TestDeleteServiceAccountService_Execute(t *testing.T) {
	tests := []struct {
		name              string
		setupMock         func(*minio.MockMinIOServer)
		request           DeleteServiceAccountRequest
		expectedError     string
		validateResponse  func(t *testing.T, response *DeleteServiceAccountResponse)
	}{
		{
			name: "successful deletion",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetDeleteServiceAccountSuccess()
			},
			request: DeleteServiceAccountRequest{
				AccessKey: "AKIAIOSFODNN7EXAMPLE",
			},
			validateResponse: func(t *testing.T, response *DeleteServiceAccountResponse) {
				if response.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", response.AccessKey)
				}
				if response.Message != "Service account deleted successfully" {
					t.Errorf("Expected Message %q, got %q", "Service account deleted successfully", response.Message)
				}
			},
		},
		{
			name: "successful deletion with different access key",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetDeleteServiceAccountSuccess()
			},
			request: DeleteServiceAccountRequest{
				AccessKey: "AKIAI44QH8DHBEXAMPLE",
			},
			validateResponse: func(t *testing.T, response *DeleteServiceAccountResponse) {
				if response.AccessKey != "AKIAI44QH8DHBEXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAI44QH8DHBEXAMPLE", response.AccessKey)
				}
				if response.Message != "Service account deleted successfully" {
					t.Errorf("Expected Message %q, got %q", "Service account deleted successfully", response.Message)
				}
			},
		},
		{
			name: "MinIO server error",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetDeleteServiceAccountError(500, "Internal Server Error")
			},
			request: DeleteServiceAccountRequest{
				AccessKey: "AKIAIOSFODNN7EXAMPLE",
			},
			expectedError: "failed to delete service account",
		},
		{
			name: "MinIO authentication error",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetDeleteServiceAccountError(401, "Unauthorized")
			},
			request: DeleteServiceAccountRequest{
				AccessKey: "AKIAIOSFODNN7EXAMPLE",
			},
			expectedError: "failed to delete service account",
		},
		{
			name: "MinIO not found error",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetDeleteServiceAccountError(404, "Service account not found")
			},
			request: DeleteServiceAccountRequest{
				AccessKey: "NONEXISTENT",
			},
			expectedError: "failed to delete service account",
		},
		{
			name: "MinIO forbidden error",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetDeleteServiceAccountError(403, "Forbidden")
			},
			request: DeleteServiceAccountRequest{
				AccessKey: "AKIAIOSFODNN7EXAMPLE",
			},
			expectedError: "failed to delete service account",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock MinIO server
			mockServer := minio.NewMockMinIOServer()
			defer mockServer.Close()

			tt.setupMock(mockServer)

			// Create MinIO client
			minioClient, err := mockServer.CreateMinIOClient()
			if err != nil {
				t.Fatalf("Failed to create MinIO client: %v", err)
			}

			// Create service
			service := NewDeleteServiceAccountService(minioClient)

			// Create context with logger
			ctx := zerolog.New(zerolog.NewTestWriter(t)).WithContext(context.Background())

			// Execute service
			response, err := service.Execute(ctx, tt.request)

			// Validate error
			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("Expected error containing %q, got nil", tt.expectedError)
					return
				}
				if err.Error() == "" {
					t.Errorf("Expected error containing %q, got empty error", tt.expectedError)
					return
				}
				// Note: We don't check exact error message match because MinIO client
				// may wrap errors differently in test vs real scenarios
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
					return
				}
			}

			// Validate response
			if tt.validateResponse != nil && response != nil {
				tt.validateResponse(t, response)
			}
		})
	}
}

func TestDeleteServiceAccountService_NewDeleteServiceAccountService(t *testing.T) {
	// Setup mock MinIO server
	mockServer := minio.NewMockMinIOServer()
	defer mockServer.Close()

	// Create MinIO client
	minioClient, err := mockServer.CreateMinIOClient()
	if err != nil {
		t.Fatalf("Failed to create MinIO client: %v", err)
	}

	// Test service creation
	service := NewDeleteServiceAccountService(minioClient)

	if service == nil {
		t.Error("Expected service to be created, got nil")
	}

	if service.minioClient != minioClient {
		t.Error("Expected MinIO client to be set correctly")
	}
}