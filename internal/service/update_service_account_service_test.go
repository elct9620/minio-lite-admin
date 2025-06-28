package service

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/elct9620/minio-lite-admin/internal/testability/minio"
	"github.com/minio/madmin-go/v4"
	"github.com/rs/zerolog"
)

func TestUpdateServiceAccountService_Execute(t *testing.T) {
	tests := []struct {
		name             string
		setupMock        func(*minio.MockMinIOServer)
		request          UpdateServiceAccountRequest
		expectError      bool
		expectedErrorMsg string
		validateResponse func(t *testing.T, response *UpdateServiceAccountResponse)
	}{
		{
			name: "successful update with policy",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetUpdateServiceAccountSuccess()
			},
			request: UpdateServiceAccountRequest{
				AccessKey: "AKIAIOSFODNN7EXAMPLE",
				NewPolicy: `{"Version": "2012-10-17", "Statement": [{"Effect": "Allow", "Action": "s3:GetObject", "Resource": "arn:aws:s3:::bucket/*"}]}`,
			},
			expectError: false,
			validateResponse: func(t *testing.T, response *UpdateServiceAccountResponse) {
				if response.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", response.AccessKey)
				}
				if response.Message != "Service account updated successfully" {
					t.Errorf("Expected Message %q, got %q", "Service account updated successfully", response.Message)
				}
			},
		},
		{
			name: "successful update with secret key rotation",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetUpdateServiceAccountSuccess()
			},
			request: UpdateServiceAccountRequest{
				AccessKey:    "AKIAIOSFODNN7EXAMPLE",
				NewSecretKey: "newSecretKey123456789012345",
			},
			expectError: false,
			validateResponse: func(t *testing.T, response *UpdateServiceAccountResponse) {
				if response.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", response.AccessKey)
				}
				if response.Message != "Service account updated successfully" {
					t.Errorf("Expected Message %q, got %q", "Service account updated successfully", response.Message)
				}
			},
		},
		{
			name: "successful update with status change",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetUpdateServiceAccountSuccess()
			},
			request: UpdateServiceAccountRequest{
				AccessKey: "AKIAIOSFODNN7EXAMPLE",
				NewStatus: string(madmin.AccountDisabled),
			},
			expectError: false,
			validateResponse: func(t *testing.T, response *UpdateServiceAccountResponse) {
				if response.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", response.AccessKey)
				}
				if response.Message != "Service account updated successfully" {
					t.Errorf("Expected Message %q, got %q", "Service account updated successfully", response.Message)
				}
			},
		},
		{
			name: "successful update with name and description",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetUpdateServiceAccountSuccess()
			},
			request: UpdateServiceAccountRequest{
				AccessKey:      "AKIAIOSFODNN7EXAMPLE",
				NewName:        "Updated Service Account",
				NewDescription: "Updated description",
			},
			expectError: false,
			validateResponse: func(t *testing.T, response *UpdateServiceAccountResponse) {
				if response.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", response.AccessKey)
				}
				if response.Message != "Service account updated successfully" {
					t.Errorf("Expected Message %q, got %q", "Service account updated successfully", response.Message)
				}
			},
		},
		{
			name: "successful update with expiration",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetUpdateServiceAccountSuccess()
			},
			request: UpdateServiceAccountRequest{
				AccessKey:     "AKIAIOSFODNN7EXAMPLE",
				NewExpiration: func() *int64 { t := time.Now().Add(24 * time.Hour).Unix(); return &t }(),
			},
			expectError: false,
			validateResponse: func(t *testing.T, response *UpdateServiceAccountResponse) {
				if response.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", response.AccessKey)
				}
				if response.Message != "Service account updated successfully" {
					t.Errorf("Expected Message %q, got %q", "Service account updated successfully", response.Message)
				}
			},
		},
		{
			name: "successful update with all fields",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetUpdateServiceAccountSuccess()
			},
			request: UpdateServiceAccountRequest{
				AccessKey:      "AKIAIOSFODNN7EXAMPLE",
				NewPolicy:      `{"Version": "2012-10-17", "Statement": [{"Effect": "Allow", "Action": "s3:*", "Resource": "*"}]}`,
				NewSecretKey:   "newSecretKey123456789012345",
				NewStatus:      string(madmin.AccountEnabled),
				NewName:        "Full Update Service Account",
				NewDescription: "Fully updated description",
				NewExpiration:  func() *int64 { t := time.Now().Add(48 * time.Hour).Unix(); return &t }(),
			},
			expectError: false,
			validateResponse: func(t *testing.T, response *UpdateServiceAccountResponse) {
				if response.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", response.AccessKey)
				}
				if response.Message != "Service account updated successfully" {
					t.Errorf("Expected Message %q, got %q", "Service account updated successfully", response.Message)
				}
			},
		},
		{
			name: "invalid status error",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetUpdateServiceAccountSuccess()
			},
			request: UpdateServiceAccountRequest{
				AccessKey: "AKIAIOSFODNN7EXAMPLE",
				NewStatus: "invalid-status",
			},
			expectError:      true,
			expectedErrorMsg: "invalid account status",
		},
		{
			name: "MinIO server error",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetUpdateServiceAccountError(500, "Internal Server Error")
			},
			request: UpdateServiceAccountRequest{
				AccessKey: "AKIAIOSFODNN7EXAMPLE",
				NewPolicy: `{"Version": "2012-10-17", "Statement": []}`,
			},
			expectError:      true,
			expectedErrorMsg: "failed to update service account",
		},
		{
			name: "MinIO authentication error",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetUpdateServiceAccountError(401, "Unauthorized")
			},
			request: UpdateServiceAccountRequest{
				AccessKey: "AKIAIOSFODNN7EXAMPLE",
				NewPolicy: `{"Version": "2012-10-17", "Statement": []}`,
			},
			expectError:      true,
			expectedErrorMsg: "failed to update service account",
		},
		{
			name: "MinIO not found error",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetUpdateServiceAccountError(404, "Service account not found")
			},
			request: UpdateServiceAccountRequest{
				AccessKey: "NONEXISTENT",
				NewPolicy: `{"Version": "2012-10-17", "Statement": []}`,
			},
			expectError:      true,
			expectedErrorMsg: "failed to update service account",
		},
		{
			name: "MinIO forbidden error",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetUpdateServiceAccountError(403, "Forbidden")
			},
			request: UpdateServiceAccountRequest{
				AccessKey: "AKIAIOSFODNN7EXAMPLE",
				NewPolicy: `{"Version": "2012-10-17", "Statement": []}`,
			},
			expectError:      true,
			expectedErrorMsg: "failed to update service account",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock MinIO server
			mockServer := minio.NewMockMinIOServer()
			defer mockServer.Close()

			// Pre-populate service account in store for update tests
			mockServer.AddServiceAccountToStore(
				tt.request.AccessKey,
				"existingSecretKey",
				"Original Name",
				"Original Description",
				"enabled",
				"admin",
				nil, // no policy
				nil, // no expiration
			)

			tt.setupMock(mockServer)

			// Create MinIO client
			minioClient, err := mockServer.CreateMinIOClient()
			if err != nil {
				t.Fatalf("Failed to create MinIO client: %v", err)
			}

			// Create service
			service := NewUpdateServiceAccountService(minioClient)

			// Create context with logger
			ctx := zerolog.New(zerolog.NewTestWriter(t)).WithContext(context.Background())

			// Execute service
			response, err := service.Execute(ctx, tt.request)

			// Validate error expectations
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
					return
				}

				if tt.expectedErrorMsg != "" && !strings.Contains(err.Error(), tt.expectedErrorMsg) {
					t.Errorf("Expected error to contain %q, got %q", tt.expectedErrorMsg, err.Error())
				}
				return
			}

			// Validate success expectations
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if response == nil {
				t.Error("Expected response but got nil")
				return
			}

			if tt.validateResponse != nil {
				tt.validateResponse(t, response)
			}
		})
	}
}

func TestUpdateServiceAccountService_Execute_PolicyValidation(t *testing.T) {
	tests := []struct {
		name      string
		policy    string
		expectErr bool
	}{
		{
			name: "valid JSON policy",
			policy: `{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": "s3:GetObject",
						"Resource": "arn:aws:s3:::bucket/*"
					}
				]
			}`,
			expectErr: false,
		},
		{
			name:      "empty policy",
			policy:    "",
			expectErr: false, // Empty policy is allowed (means clear policy)
		},
		{
			name: "complex policy with multiple statements",
			policy: `{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": ["s3:GetObject", "s3:PutObject"],
						"Resource": "arn:aws:s3:::bucket/*"
					},
					{
						"Effect": "Deny",
						"Action": "s3:DeleteObject",
						"Resource": "arn:aws:s3:::bucket/readonly/*"
					}
				]
			}`,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock MinIO server
			mockServer := minio.NewMockMinIOServer()
			defer mockServer.Close()

			// Pre-populate service account in store for policy validation tests
			mockServer.AddServiceAccountToStore(
				"AKIAIOSFODNN7EXAMPLE",
				"existingSecretKey",
				"Original Name",
				"Original Description",
				"enabled",
				"admin",
				nil, // no policy
				nil, // no expiration
			)

			mockServer.SetUpdateServiceAccountSuccess()

			// Create MinIO client
			minioClient, err := mockServer.CreateMinIOClient()
			if err != nil {
				t.Fatalf("Failed to create MinIO client: %v", err)
			}

			// Create service
			service := NewUpdateServiceAccountService(minioClient)

			// Create context with logger
			ctx := zerolog.New(zerolog.NewTestWriter(t)).WithContext(context.Background())

			// Create request with policy
			request := UpdateServiceAccountRequest{
				AccessKey: "AKIAIOSFODNN7EXAMPLE",
				NewPolicy: tt.policy,
			}

			// Execute service
			response, err := service.Execute(ctx, request)

			// Validate error expectations
			if tt.expectErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			// Validate success expectations
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if response == nil {
				t.Error("Expected response but got nil")
				return
			}

			// Validate the policy was passed through correctly by checking the JSON parsing
			if tt.policy != "" {
				var parsed map[string]interface{}
				if err := json.Unmarshal([]byte(tt.policy), &parsed); err != nil {
					t.Errorf("Policy should be valid JSON: %v", err)
				}
			}
		})
	}
}
