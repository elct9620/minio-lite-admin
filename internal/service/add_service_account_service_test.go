package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/elct9620/minio-lite-admin/internal/testability/minio"
	"github.com/rs/zerolog"
)

func TestAddServiceAccountService_Execute(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*minio.MockMinIOServer)
		request        CreateServiceAccountRequest
		expectedError  string
		validateResult func(t *testing.T, result *CreateServiceAccountResponse)
	}{
		{
			name: "successful creation with auto-generated keys",
			setupMock: func(mock *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				mock.SetAddServiceAccountResponse(scenarios.SuccessfulAddServiceAccount())
			},
			request: CreateServiceAccountRequest{
				Name:        "test-service-account",
				Description: "Test service account",
			},
			validateResult: func(t *testing.T, result *CreateServiceAccountResponse) {
				if result.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", result.AccessKey)
				}
				if result.SecretKey != "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" {
					t.Errorf("Expected SecretKey %q, got %q", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", result.SecretKey)
				}
				if result.Name != "test-service-account" {
					t.Errorf("Expected Name %q, got %q", "test-service-account", result.Name)
				}
				if result.Description != "Test service account" {
					t.Errorf("Expected Description %q, got %q", "Test service account", result.Description)
				}
				if result.SessionToken != "" {
					t.Errorf("Expected empty SessionToken, got %q", result.SessionToken)
				}
				if !result.Expiration.IsZero() {
					t.Errorf("Expected zero Expiration, got %v", result.Expiration)
				}
			},
		},
		{
			name: "successful creation with custom keys",
			setupMock: func(mock *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				mock.SetAddServiceAccountResponse(scenarios.CustomKeysAddServiceAccount("CUSTOM123", "customSecret456"))
			},
			request: CreateServiceAccountRequest{
				Name:        "custom-service-account",
				Description: "Custom service account",
				AccessKey:   "CUSTOM123",
				SecretKey:   "customSecret456",
			},
			validateResult: func(t *testing.T, result *CreateServiceAccountResponse) {
				if result.AccessKey != "CUSTOM123" {
					t.Errorf("Expected AccessKey %q, got %q", "CUSTOM123", result.AccessKey)
				}
				if result.SecretKey != "customSecret456" {
					t.Errorf("Expected SecretKey %q, got %q", "customSecret456", result.SecretKey)
				}
				if result.Name != "custom-service-account" {
					t.Errorf("Expected Name %q, got %q", "custom-service-account", result.Name)
				}
				if result.Description != "Custom service account" {
					t.Errorf("Expected Description %q, got %q", "Custom service account", result.Description)
				}
			},
		},
		{
			name: "successful creation with policy",
			setupMock: func(mock *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				mock.SetAddServiceAccountResponse(scenarios.SuccessfulAddServiceAccount())
			},
			request: CreateServiceAccountRequest{
				Name:   "policy-service-account",
				Policy: `{"Version": "2012-10-17", "Statement": [{"Effect": "Allow", "Action": ["s3:GetObject"], "Resource": ["arn:aws:s3:::bucket/*"]}]}`,
			},
			validateResult: func(t *testing.T, result *CreateServiceAccountResponse) {
				if result.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", result.AccessKey)
				}
				if result.Name != "policy-service-account" {
					t.Errorf("Expected Name %q, got %q", "policy-service-account", result.Name)
				}
			},
		},
		{
			name: "successful creation with expiration",
			setupMock: func(mock *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				expiration := time.Now().Add(24 * time.Hour)
				mock.SetAddServiceAccountResponse(scenarios.ExpiringAddServiceAccount(expiration))
			},
			request: CreateServiceAccountRequest{
				Name: "expiring-service-account",
				Expiration: func() *string {
					exp := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
					return &exp
				}(),
			},
			validateResult: func(t *testing.T, result *CreateServiceAccountResponse) {
				if result.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", result.AccessKey)
				}
				if result.Name != "expiring-service-account" {
					t.Errorf("Expected Name %q, got %q", "expiring-service-account", result.Name)
				}
				if result.Expiration.IsZero() {
					t.Errorf("Expected non-zero Expiration, got zero")
				}
			},
		},
		{
			name: "successful creation with target user",
			setupMock: func(mock *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				mock.SetAddServiceAccountResponse(scenarios.SuccessfulAddServiceAccount())
			},
			request: CreateServiceAccountRequest{
				Name:       "target-user-service-account",
				TargetUser: "testuser",
			},
			validateResult: func(t *testing.T, result *CreateServiceAccountResponse) {
				if result.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", result.AccessKey)
				}
				if result.Name != "target-user-service-account" {
					t.Errorf("Expected Name %q, got %q", "target-user-service-account", result.Name)
				}
			},
		},
		{
			name: "successful creation with all fields",
			setupMock: func(mock *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				expiration := time.Now().Add(24 * time.Hour)
				mock.SetAddServiceAccountResponse(scenarios.ExpiringAddServiceAccount(expiration))
			},
			request: CreateServiceAccountRequest{
				Name:        "complete-service-account",
				Description: "Complete service account with all fields",
				AccessKey:   "COMPLETE123",
				SecretKey:   "completeSecret456",
				Policy:      `{"Version": "2012-10-17", "Statement": [{"Effect": "Allow", "Action": ["s3:*"], "Resource": ["*"]}]}`,
				TargetUser:  "admin",
				Expiration: func() *string {
					exp := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
					return &exp
				}(),
			},
			validateResult: func(t *testing.T, result *CreateServiceAccountResponse) {
				if result.AccessKey != "AKIAIOSFODNN7EXAMPLE" {
					t.Errorf("Expected AccessKey %q, got %q", "AKIAIOSFODNN7EXAMPLE", result.AccessKey)
				}
				if result.Name != "complete-service-account" {
					t.Errorf("Expected Name %q, got %q", "complete-service-account", result.Name)
				}
				if result.Description != "Complete service account with all fields" {
					t.Errorf("Expected Description %q, got %q", "Complete service account with all fields", result.Description)
				}
			},
		},
		{
			name:      "invalid expiration format",
			setupMock: func(mock *minio.MockMinIOServer) {},
			request: CreateServiceAccountRequest{
				Name:       "invalid-expiration",
				Expiration: func() *string { exp := "invalid-date"; return &exp }(),
			},
			expectedError: "invalid expiration format, expected RFC3339",
		},
		{
			name: "MinIO server error",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetAddServiceAccountError(500, "Internal Server Error")
			},
			request: CreateServiceAccountRequest{
				Name: "error-service-account",
			},
			expectedError: "failed to create service account",
		},
		{
			name: "MinIO authentication error",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetAddServiceAccountError(401, "Unauthorized")
			},
			request: CreateServiceAccountRequest{
				Name: "auth-error-service-account",
			},
			expectedError: "failed to create service account",
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
			service := NewAddServiceAccountService(minioClient)

			// Create context with logger
			ctx := zerolog.New(zerolog.NewTestWriter(t)).WithContext(context.Background())

			// Execute test
			result, err := service.Execute(ctx, tt.request)

			// Validate results
			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("Expected error containing %q, but got no error", tt.expectedError)
				} else if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("Expected error containing %q, got %q", tt.expectedError, err.Error())
				}
				if result != nil {
					t.Errorf("Expected nil result when error occurs, got %+v", result)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if result == nil {
					t.Fatal("Expected result, got nil")
				}
				if tt.validateResult != nil {
					tt.validateResult(t, result)
				}
			}
		})
	}
}

// TestAddServiceAccountService_NewAddServiceAccountService tests the constructor
func TestAddServiceAccountService_NewAddServiceAccountService(t *testing.T) {
	mockServer := minio.NewMockMinIOServer()
	defer mockServer.Close()

	minioClient, err := mockServer.CreateMinIOClient()
	if err != nil {
		t.Fatalf("Failed to create MinIO client: %v", err)
	}

	service := NewAddServiceAccountService(minioClient)
	if service == nil {
		t.Fatal("Expected service, got nil")
	}
	if service.minioClient != minioClient {
		t.Errorf("Expected service.minioClient to be %v, got %v", minioClient, service.minioClient)
	}
}
