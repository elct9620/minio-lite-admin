package service

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/madmin-go/v4"
	"github.com/rs/zerolog"
)

type AddServiceAccountService struct {
	minioClient *madmin.AdminClient
}

// CreateServiceAccountRequest represents the request to create a service account
type CreateServiceAccountRequest struct {
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	AccessKey   string  `json:"accessKey,omitempty"`   // Optional - MinIO generates if empty
	SecretKey   string  `json:"secretKey,omitempty"`   // Optional - MinIO generates if empty
	Policy      string  `json:"policy,omitempty"`      // JSON policy document
	TargetUser  string  `json:"targetUser,omitempty"`  // User this service account belongs to
	Expiration  *string `json:"expiration,omitempty"`  // ISO 8601 format
}

// CreateServiceAccountResponse represents the response from creating a service account
type CreateServiceAccountResponse struct {
	AccessKey    string    `json:"accessKey"`
	SecretKey    string    `json:"secretKey"`
	SessionToken string    `json:"sessionToken,omitempty"`
	Expiration   time.Time `json:"expiration,omitempty"`
	Name         string    `json:"name,omitempty"`
	Description  string    `json:"description,omitempty"`
}

func NewAddServiceAccountService(minioClient *madmin.AdminClient) *AddServiceAccountService {
	return &AddServiceAccountService{
		minioClient: minioClient,
	}
}

func (s *AddServiceAccountService) Execute(ctx context.Context, req CreateServiceAccountRequest) (*CreateServiceAccountResponse, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().
		Str("name", req.Name).
		Str("targetUser", req.TargetUser).
		Bool("hasPolicy", req.Policy != "").
		Bool("hasCustomAccessKey", req.AccessKey != "").
		Msg("Creating service account")

	// Prepare the MinIO request
	addReq := madmin.AddServiceAccountReq{
		Name:        req.Name,
		Description: req.Description,
		AccessKey:   req.AccessKey,
		SecretKey:   req.SecretKey,
		TargetUser:  req.TargetUser,
	}

	// Parse policy if provided
	if req.Policy != "" {
		addReq.Policy = []byte(req.Policy)
	}

	// Parse expiration if provided
	if req.Expiration != nil && *req.Expiration != "" {
		expTime, err := time.Parse(time.RFC3339, *req.Expiration)
		if err != nil {
			logger.Error().Err(err).Str("expiration", *req.Expiration).Msg("Failed to parse expiration time")
			return nil, fmt.Errorf("invalid expiration format, expected RFC3339: %w", err)
		}
		addReq.Expiration = &expTime
	}

	// Create the service account
	credentials, err := s.minioClient.AddServiceAccount(ctx, addReq)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create service account")
		return nil, fmt.Errorf("failed to create service account: %w", err)
	}

	logger.Info().
		Str("accessKey", credentials.AccessKey).
		Str("name", req.Name).
		Msg("Successfully created service account")

	// Prepare response
	response := &CreateServiceAccountResponse{
		AccessKey:    credentials.AccessKey,
		SecretKey:    credentials.SecretKey,
		SessionToken: credentials.SessionToken,
		Expiration:   credentials.Expiration,
		Name:         req.Name,
		Description:  req.Description,
	}

	return response, nil
}