package service

import (
	"context"
	"fmt"

	"github.com/minio/madmin-go/v4"
	"github.com/rs/zerolog"
)

type ListAccessKeysService struct {
	minioClient *madmin.AdminClient
}

// AccessKeyInfo represents a unified access key information structure
type AccessKeyInfo struct {
	AccessKey     string  `json:"accessKey"`
	ParentUser    string  `json:"parentUser"`
	AccountStatus string  `json:"accountStatus"`
	Type          string  `json:"type"` // "user", "serviceAccount", "sts"
	Name          string  `json:"name,omitempty"`
	Description   string  `json:"description,omitempty"`
	Expiration    *string `json:"expiration,omitempty"` // ISO 8601 format
	CreatedAt     *string `json:"createdAt,omitempty"`  // ISO 8601 format
	ImpliedPolicy bool    `json:"impliedPolicy"`
}

// ListAccessKeysResponse represents the API response for listing access keys
type ListAccessKeysResponse struct {
	AccessKeys []AccessKeyInfo `json:"accessKeys"`
	Total      int             `json:"total"`
}

// ListAccessKeysOptions represents options for filtering access keys
type ListAccessKeysOptions struct {
	Type string // "all", "users", "serviceAccounts", "sts"
	User string // Filter by specific user (optional)
}

func NewListAccessKeysService(minioClient *madmin.AdminClient) *ListAccessKeysService {
	return &ListAccessKeysService{
		minioClient: minioClient,
	}
}

func (s *ListAccessKeysService) Execute(ctx context.Context, opts ListAccessKeysOptions) (*ListAccessKeysResponse, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Str("type", opts.Type).Str("user", opts.User).Msg("Listing access keys")

	var allAccessKeys []AccessKeyInfo

	// Get all users first to get the list of users for bulk access key retrieval
	users, err := s.minioClient.ListUsers(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to list users")
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Note: We don't need to extract userNames as we use ListAccessKeysBulk with opts.User directly

	// Use ListAccessKeysBulk to get all access keys for users
	bulkOpts := madmin.ListAccessKeysOpts{}

	// Set list type based on filter
	switch opts.Type {
	case "users":
		bulkOpts.ListType = madmin.AccessKeyListUsersOnly
	case "serviceAccounts":
		bulkOpts.ListType = madmin.AccessKeyListSvcaccOnly
	case "sts":
		bulkOpts.ListType = madmin.AccessKeyListSTSOnly
	default:
		bulkOpts.ListType = madmin.AccessKeyListAll
	}

	// If no specific user filter, get all users
	var accessKeysMap map[string]madmin.ListAccessKeysResp

	if opts.User == "" {
		// Use All: true when no specific user is requested
		bulkOpts.All = true
		accessKeysMap, err = s.minioClient.ListAccessKeysBulk(ctx, nil, bulkOpts)
	} else {
		// Use specific user list when filtering by user
		accessKeysMap, err = s.minioClient.ListAccessKeysBulk(ctx, []string{opts.User}, bulkOpts)
	}
	if err != nil {
		logger.Error().Err(err).Msg("Failed to list access keys")
		return nil, fmt.Errorf("failed to list access keys: %w", err)
	}

	// Process the results and convert to unified format
	for userName, accessKeysResp := range accessKeysMap {
		// Add user access key if it exists
		userInfo, exists := users[userName]
		if exists && (opts.Type == "all" || opts.Type == "users") {
			allAccessKeys = append(allAccessKeys, AccessKeyInfo{
				AccessKey:     userName, // For users, access key is the username
				ParentUser:    userName,
				AccountStatus: string(userInfo.Status),
				Type:          "user",
				ImpliedPolicy: false,
			})
		}

		// Add service accounts
		if opts.Type == "all" || opts.Type == "serviceAccounts" {
			for _, svcAccount := range accessKeysResp.ServiceAccounts {
				accessKey := AccessKeyInfo{
					AccessKey:     svcAccount.AccessKey,
					ParentUser:    svcAccount.ParentUser,
					AccountStatus: svcAccount.AccountStatus,
					Type:          "serviceAccount",
					Name:          svcAccount.Name,
					Description:   svcAccount.Description,
					ImpliedPolicy: svcAccount.ImpliedPolicy,
				}

				// Convert expiration to ISO 8601 string if present
				if svcAccount.Expiration != nil {
					expStr := svcAccount.Expiration.Format("2006-01-02T15:04:05Z07:00")
					accessKey.Expiration = &expStr
				}

				allAccessKeys = append(allAccessKeys, accessKey)
			}
		}

		// Add STS keys
		if opts.Type == "all" || opts.Type == "sts" {
			for _, stsKey := range accessKeysResp.STSKeys {
				accessKey := AccessKeyInfo{
					AccessKey:     stsKey.AccessKey,
					ParentUser:    stsKey.ParentUser,
					AccountStatus: stsKey.AccountStatus,
					Type:          "sts",
					Name:          stsKey.Name,
					Description:   stsKey.Description,
					ImpliedPolicy: stsKey.ImpliedPolicy,
				}

				// Convert expiration to ISO 8601 string if present
				if stsKey.Expiration != nil {
					expStr := stsKey.Expiration.Format("2006-01-02T15:04:05Z07:00")
					accessKey.Expiration = &expStr
				}

				allAccessKeys = append(allAccessKeys, accessKey)
			}
		}
	}

	logger.Debug().Int("count", len(allAccessKeys)).Msg("Successfully listed access keys")

	// Ensure accessKeys is never nil for JSON serialization
	if allAccessKeys == nil {
		allAccessKeys = []AccessKeyInfo{}
	}

	return &ListAccessKeysResponse{
		AccessKeys: allAccessKeys,
		Total:      len(allAccessKeys),
	}, nil
}
