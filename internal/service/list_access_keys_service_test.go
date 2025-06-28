package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/elct9620/minio-lite-admin/internal/testability/minio"
	"github.com/rs/zerolog"
)

func TestListAccessKeysService_Execute(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*minio.MockMinIOServer)
		options        ListAccessKeysOptions
		expectedError  string
		validateResult func(t *testing.T, result *ListAccessKeysResponse)
	}{
		{
			name: "successful list all access keys with detailed service account info",
			setupMock: func(mock *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				users, accessKeys := scenarios.SuccessfulAccessKeys()
				mock.SetUsersResponse(users)
				mock.SetAccessKeysBulkResponse(accessKeys)

				// Set detailed info for service accounts
				mock.SetInfoServiceAccountResponse("AKIAIOSFODNN7EXAMPLE", minio.InfoServiceAccountResponse{
					ParentUser:    "minioadmin",
					AccountStatus: "enabled",
					ImpliedPolicy: true,
					Name:          "enhanced-admin-account",
					Description:   "Enhanced admin service account with detailed info",
					Expiration:    nil,
				})

				mock.SetInfoServiceAccountResponse("AKIAI44QH8DHBEXAMPLE", minio.InfoServiceAccountResponse{
					ParentUser:    "minioadmin",
					AccountStatus: "enabled",
					ImpliedPolicy: false,
					Name:          "enhanced-backup-account",
					Description:   "Enhanced backup service account with detailed info",
					Expiration:    nil,
				})

				mock.SetInfoServiceAccountResponse("AKIATEST123456789012", minio.InfoServiceAccountResponse{
					ParentUser:    "testuser",
					AccountStatus: "enabled",
					ImpliedPolicy: false,
					Name:          "enhanced-test-account",
					Description:   "Enhanced test user service account with detailed info",
					Expiration:    nil,
				})
			},
			options: ListAccessKeysOptions{Type: "all"},
			validateResult: func(t *testing.T, result *ListAccessKeysResponse) {
				if result.Total != 7 { // 3 users + 3 service accounts + 1 STS key
					t.Errorf("Expected Total %d, got %d", 7, result.Total)
				}
				if len(result.AccessKeys) != 7 {
					t.Errorf("Expected %d access keys, got %d", 7, len(result.AccessKeys))
				}

				// Find service accounts and verify enhanced info
				var adminSA, backupSA, testSA *AccessKeyInfo
				for i := range result.AccessKeys {
					ak := &result.AccessKeys[i]
					if ak.AccessKey == "AKIAIOSFODNN7EXAMPLE" {
						adminSA = ak
					} else if ak.AccessKey == "AKIAI44QH8DHBEXAMPLE" {
						backupSA = ak
					} else if ak.AccessKey == "AKIATEST123456789012" {
						testSA = ak
					}
				}

				if adminSA == nil {
					t.Fatal("Expected to find admin service account")
				}
				if adminSA.Name != "enhanced-admin-account" {
					t.Errorf("Expected Name %q, got %q", "enhanced-admin-account", adminSA.Name)
				}
				if adminSA.Description != "Enhanced admin service account with detailed info" {
					t.Errorf("Expected Description %q, got %q", "Enhanced admin service account with detailed info", adminSA.Description)
				}

				if backupSA == nil {
					t.Fatal("Expected to find backup service account")
				}
				if backupSA.Name != "enhanced-backup-account" {
					t.Errorf("Expected Name %q, got %q", "enhanced-backup-account", backupSA.Name)
				}
				if backupSA.ImpliedPolicy != false {
					t.Errorf("Expected ImpliedPolicy false, got %v", backupSA.ImpliedPolicy)
				}

				if testSA == nil {
					t.Fatal("Expected to find test service account")
				}
				if testSA.Name != "enhanced-test-account" {
					t.Errorf("Expected Name %q, got %q", "enhanced-test-account", testSA.Name)
				}
				if testSA.ParentUser != "testuser" {
					t.Errorf("Expected ParentUser %q, got %q", "testuser", testSA.ParentUser)
				}
			},
		},
		{
			name: "service account info fallback when InfoServiceAccount fails",
			setupMock: func(mock *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				users, accessKeys := scenarios.SuccessfulAccessKeys()
				mock.SetUsersResponse(users)
				mock.SetAccessKeysBulkResponse(accessKeys)

				// Set error for InfoServiceAccount calls
				mock.SetInfoServiceAccountError(404, "Service account not found")
			},
			options: ListAccessKeysOptions{Type: "serviceAccounts"},
			validateResult: func(t *testing.T, result *ListAccessKeysResponse) {
				if result.Total != 3 { // 3 service accounts from bulk response
					t.Errorf("Expected Total %d, got %d", 3, result.Total)
				}
				if len(result.AccessKeys) != 3 {
					t.Errorf("Expected %d access keys, got %d", 3, len(result.AccessKeys))
				}

				// Verify fallback to basic info from ListAccessKeysBulk
				for _, ak := range result.AccessKeys {
					if ak.Type != "serviceAccount" {
						t.Errorf("Expected Type %q, got %q", "serviceAccount", ak.Type)
					}
					// Should have basic names from bulk response
					if ak.AccessKey == "AKIAIOSFODNN7EXAMPLE" && ak.Name != "admin-service-account" {
						t.Errorf("Expected fallback Name %q, got %q", "admin-service-account", ak.Name)
					}
				}
			},
		},
		{
			name: "successful list users only",
			setupMock: func(mock *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				users, accessKeys := scenarios.SuccessfulAccessKeys()
				mock.SetUsersResponse(users)
				mock.SetAccessKeysBulkResponse(accessKeys)
			},
			options: ListAccessKeysOptions{Type: "users"},
			validateResult: func(t *testing.T, result *ListAccessKeysResponse) {
				if result.Total != 3 { // 3 users
					t.Errorf("Expected Total %d, got %d", 3, result.Total)
				}
				if len(result.AccessKeys) != 3 {
					t.Errorf("Expected %d access keys, got %d", 3, len(result.AccessKeys))
				}
				for _, ak := range result.AccessKeys {
					if ak.Type != "user" {
						t.Errorf("Expected Type %q, got %q", "user", ak.Type)
					}
				}
			},
		},
		{
			name: "successful list service accounts only with enhanced info",
			setupMock: func(mock *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				users, accessKeys := scenarios.SuccessfulAccessKeys()
				mock.SetUsersResponse(users)
				mock.SetAccessKeysBulkResponse(accessKeys)

				// Add detailed info for service accounts
				expiration := time.Now().Add(30 * 24 * time.Hour) // 30 days
				mock.SetInfoServiceAccountResponse("AKIAIOSFODNN7EXAMPLE", minio.InfoServiceAccountResponse{
					ParentUser:    "minioadmin",
					AccountStatus: "enabled",
					ImpliedPolicy: true,
					Name:          "production-admin-sa",
					Description:   "Production admin service account",
					Expiration:    &expiration,
				})
			},
			options: ListAccessKeysOptions{Type: "serviceAccounts"},
			validateResult: func(t *testing.T, result *ListAccessKeysResponse) {
				if result.Total != 3 { // 3 service accounts
					t.Errorf("Expected Total %d, got %d", 3, result.Total)
				}
				if len(result.AccessKeys) != 3 {
					t.Errorf("Expected %d access keys, got %d", 3, len(result.AccessKeys))
				}
				for _, ak := range result.AccessKeys {
					if ak.Type != "serviceAccount" {
						t.Errorf("Expected Type %q, got %q", "serviceAccount", ak.Type)
					}
					// Check if enhanced info was applied
					if ak.AccessKey == "AKIAIOSFODNN7EXAMPLE" {
						if ak.Name != "production-admin-sa" {
							t.Errorf("Expected enhanced Name %q, got %q", "production-admin-sa", ak.Name)
						}
						if ak.Description != "Production admin service account" {
							t.Errorf("Expected enhanced Description %q, got %q", "Production admin service account", ak.Description)
						}
						if ak.Expiration == nil {
							t.Error("Expected expiration to be set")
						}
					}
				}
			},
		},
		{
			name: "successful list STS keys only",
			setupMock: func(mock *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				users, accessKeys := scenarios.SuccessfulAccessKeys()
				mock.SetUsersResponse(users)
				mock.SetAccessKeysBulkResponse(accessKeys)
			},
			options: ListAccessKeysOptions{Type: "sts"},
			validateResult: func(t *testing.T, result *ListAccessKeysResponse) {
				if result.Total != 1 { // 1 STS key
					t.Errorf("Expected Total %d, got %d", 1, result.Total)
				}
				if len(result.AccessKeys) != 1 {
					t.Errorf("Expected %d access keys, got %d", 1, len(result.AccessKeys))
				}
				for _, ak := range result.AccessKeys {
					if ak.Type != "sts" {
						t.Errorf("Expected Type %q, got %q", "sts", ak.Type)
					}
				}
			},
		},
		{
			name: "successful list with user filter",
			setupMock: func(mock *minio.MockMinIOServer) {
				// When filtering by specific user, only return that user's data
				users := minio.AccessKeysUsersResponse{
					"testuser": {Status: "enabled"},
				}
				accessKeys := minio.AccessKeysBulkResponse{
					"testuser": {
						ServiceAccounts: []minio.AccessKeysServiceAccount{
							{
								ParentUser:    "testuser",
								AccountStatus: "enabled",
								AccessKey:     "AKIATEST123456789012",
								Name:          "test-account",
								Description:   "Test user service account",
								ImpliedPolicy: false,
							},
						},
						STSKeys: []minio.AccessKeysServiceAccount{},
					},
				}
				mock.SetUsersResponse(users)
				mock.SetAccessKeysBulkResponse(accessKeys)
			},
			options: ListAccessKeysOptions{Type: "all", User: "testuser"},
			validateResult: func(t *testing.T, result *ListAccessKeysResponse) {
				// Should return testuser + their service account
				if result.Total != 2 {
					t.Errorf("Expected Total %d, got %d", 2, result.Total)
				}
				for _, ak := range result.AccessKeys {
					if ak.ParentUser != "testuser" && ak.AccessKey != "testuser" {
						t.Errorf("Expected access key related to testuser, got %+v", ak)
					}
				}
			},
		},
		{
			name: "empty response when no access keys exist",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetUsersResponse(minio.AccessKeysUsersResponse{})
				mock.SetAccessKeysBulkResponse(minio.AccessKeysBulkResponse{})
			},
			options: ListAccessKeysOptions{Type: "all"},
			validateResult: func(t *testing.T, result *ListAccessKeysResponse) {
				if result.Total != 0 {
					t.Errorf("Expected Total %d, got %d", 0, result.Total)
				}
				if len(result.AccessKeys) != 0 {
					t.Errorf("Expected %d access keys, got %d", 0, len(result.AccessKeys))
				}
			},
		},
		{
			name: "MinIO list users error",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetUsersError(500, "Internal Server Error")
			},
			options:       ListAccessKeysOptions{Type: "all"},
			expectedError: "failed to list users",
		},
		{
			name: "MinIO list access keys bulk error",
			setupMock: func(mock *minio.MockMinIOServer) {
				scenarios := minio.TestScenarios{}
				users, _ := scenarios.SuccessfulAccessKeys()
				mock.SetUsersResponse(users)
				mock.SetAccessKeysBulkError(500, "Internal Server Error")
			},
			options:       ListAccessKeysOptions{Type: "all"},
			expectedError: "failed to list access keys",
		},
		{
			name: "MinIO authentication error",
			setupMock: func(mock *minio.MockMinIOServer) {
				mock.SetUsersError(401, "Unauthorized")
			},
			options:       ListAccessKeysOptions{Type: "all"},
			expectedError: "failed to list users",
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
			service := NewListAccessKeysService(minioClient)

			// Create context with logger
			ctx := zerolog.New(zerolog.NewTestWriter(t)).WithContext(context.Background())

			// Execute test
			result, err := service.Execute(ctx, tt.options)

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

// TestListAccessKeysService_NewListAccessKeysService tests the constructor
func TestListAccessKeysService_NewListAccessKeysService(t *testing.T) {
	mockServer := minio.NewMockMinIOServer()
	defer mockServer.Close()

	minioClient, err := mockServer.CreateMinIOClient()
	if err != nil {
		t.Fatalf("Failed to create MinIO client: %v", err)
	}

	service := NewListAccessKeysService(minioClient)
	if service == nil {
		t.Fatal("Expected service, got nil")
	}
	if service.minioClient != minioClient {
		t.Errorf("Expected service.minioClient to be %v, got %v", minioClient, service.minioClient)
	}
}

// TestListAccessKeysService_InfoServiceAccountIntegration tests the InfoServiceAccount integration specifically
func TestListAccessKeysService_InfoServiceAccountIntegration(t *testing.T) {
	tests := []struct {
		name             string
		setupMock        func(*minio.MockMinIOServer)
		expectedEnhanced bool
		expectedName     string
		expectedDesc     string
		accessKeyToCheck string
	}{
		{
			name: "InfoServiceAccount success enhances basic info",
			setupMock: func(mock *minio.MockMinIOServer) {
				// Basic response from ListAccessKeysBulk
				users := minio.AccessKeysUsersResponse{
					"minioadmin": {Status: "enabled"},
				}
				accessKeys := minio.AccessKeysBulkResponse{
					"minioadmin": {
						ServiceAccounts: []minio.AccessKeysServiceAccount{
							{
								ParentUser:    "minioadmin",
								AccountStatus: "enabled",
								AccessKey:     "BASIC123",
								Name:          "basic-name",
								Description:   "basic-description",
								ImpliedPolicy: false,
							},
						},
					},
				}
				mock.SetUsersResponse(users)
				mock.SetAccessKeysBulkResponse(accessKeys)

				// Enhanced response from InfoServiceAccount
				mock.SetInfoServiceAccountResponse("BASIC123", minio.InfoServiceAccountResponse{
					ParentUser:    "minioadmin",
					AccountStatus: "enabled",
					ImpliedPolicy: true,                   // Different from basic
					Name:          "enhanced-name",        // Different from basic
					Description:   "enhanced-description", // Different from basic
				})
			},
			expectedEnhanced: true,
			expectedName:     "enhanced-name",
			expectedDesc:     "enhanced-description",
			accessKeyToCheck: "BASIC123",
		},
		{
			name: "InfoServiceAccount failure falls back to basic info",
			setupMock: func(mock *minio.MockMinIOServer) {
				// Basic response from ListAccessKeysBulk
				users := minio.AccessKeysUsersResponse{
					"minioadmin": {Status: "enabled"},
				}
				accessKeys := minio.AccessKeysBulkResponse{
					"minioadmin": {
						ServiceAccounts: []minio.AccessKeysServiceAccount{
							{
								ParentUser:    "minioadmin",
								AccountStatus: "enabled",
								AccessKey:     "FALLBACK123",
								Name:          "fallback-name",
								Description:   "fallback-description",
								ImpliedPolicy: false,
							},
						},
					},
				}
				mock.SetUsersResponse(users)
				mock.SetAccessKeysBulkResponse(accessKeys)

				// Error response from InfoServiceAccount
				mock.SetInfoServiceAccountError(404, "Service account not found")
			},
			expectedEnhanced: false,
			expectedName:     "fallback-name",
			expectedDesc:     "fallback-description",
			accessKeyToCheck: "FALLBACK123",
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
			service := NewListAccessKeysService(minioClient)

			// Create context with logger
			ctx := zerolog.New(zerolog.NewTestWriter(t)).WithContext(context.Background())

			// Execute test
			result, err := service.Execute(ctx, ListAccessKeysOptions{Type: "serviceAccounts"})

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result == nil {
				t.Fatal("Expected result, got nil")
			}

			// Find the service account we're checking
			var foundSA *AccessKeyInfo
			for i := range result.AccessKeys {
				if result.AccessKeys[i].AccessKey == tt.accessKeyToCheck {
					foundSA = &result.AccessKeys[i]
					break
				}
			}

			if foundSA == nil {
				t.Fatalf("Expected to find service account %q", tt.accessKeyToCheck)
			}

			// Validate enhanced/fallback behavior
			if foundSA.Name != tt.expectedName {
				t.Errorf("Expected Name %q, got %q", tt.expectedName, foundSA.Name)
			}
			if foundSA.Description != tt.expectedDesc {
				t.Errorf("Expected Description %q, got %q", tt.expectedDesc, foundSA.Description)
			}

			// Check ImpliedPolicy to verify enhanced vs fallback
			if tt.expectedEnhanced && !foundSA.ImpliedPolicy {
				t.Error("Expected enhanced info to have ImpliedPolicy=true")
			}
			if !tt.expectedEnhanced && foundSA.ImpliedPolicy {
				t.Error("Expected fallback info to have ImpliedPolicy=false")
			}
		})
	}
}
