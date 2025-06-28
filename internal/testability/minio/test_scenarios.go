package minio

import "time"

// TestScenarios provides pre-configured test scenarios
type TestScenarios struct{}

// Server Info Scenarios

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

// Access Keys Scenarios

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

// Service Account Creation Scenarios

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

// Service Account Update Scenarios

// SuccessfulUpdateServiceAccount returns a successful update operation
func (TestScenarios) SuccessfulUpdateServiceAccount() bool {
	return true
}

// Service Account Deletion Scenarios

// SuccessfulDeleteServiceAccount returns a successful delete operation
func (TestScenarios) SuccessfulDeleteServiceAccount() bool {
	return true
}

// Service Account Info Scenarios

// SuccessfulInfoServiceAccount returns a typical service account info response
func (TestScenarios) SuccessfulInfoServiceAccount(accessKey string) InfoServiceAccountResponse {
	return InfoServiceAccountResponse{
		ParentUser:    "minioadmin",
		AccountStatus: "enabled",
		ImpliedPolicy: true,
		Name:          "detailed-service-account",
		Description:   "Service account with detailed information",
		Expiration:    nil,
	}
}
