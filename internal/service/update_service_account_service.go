package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/minio/madmin-go/v4"
	"github.com/rs/zerolog"
)

type UpdateServiceAccountService struct {
	minioClient *madmin.AdminClient
}

type UpdateServiceAccountRequest struct {
	AccessKey      string `json:"accessKey"`
	NewPolicy      string `json:"newPolicy,omitempty"`
	NewSecretKey   string `json:"newSecretKey,omitempty"`
	NewStatus      string `json:"newStatus,omitempty"`
	NewName        string `json:"newName,omitempty"`
	NewDescription string `json:"newDescription,omitempty"`
	NewExpiration  *int64 `json:"newExpiration,omitempty"` // Unix timestamp in seconds
}

type UpdateServiceAccountResponse struct {
	AccessKey    string `json:"accessKey"`
	Message      string `json:"message"`
	NewSecretKey string `json:"newSecretKey,omitempty"` // Only populated when secret key is rotated
}

func NewUpdateServiceAccountService(minioClient *madmin.AdminClient) *UpdateServiceAccountService {
	return &UpdateServiceAccountService{
		minioClient: minioClient,
	}
}

func (s *UpdateServiceAccountService) Execute(ctx context.Context, req UpdateServiceAccountRequest) (*UpdateServiceAccountResponse, error) {
	logger := zerolog.Ctx(ctx)

	logger.Debug().
		Str("accessKey", req.AccessKey).
		Bool("hasNewPolicy", req.NewPolicy != "").
		Bool("hasNewSecretKey", req.NewSecretKey != "").
		Str("newStatus", req.NewStatus).
		Str("newName", req.NewName).
		Str("newDescription", req.NewDescription).
		Msg("Executing update service account request")

	// Prepare the update request
	updateReq := madmin.UpdateServiceAccountReq{}

	// Set policy if provided
	if req.NewPolicy != "" {
		updateReq.NewPolicy = json.RawMessage(req.NewPolicy)
	}

	// Set secret key if provided
	if req.NewSecretKey != "" {
		updateReq.NewSecretKey = req.NewSecretKey
	}

	// Set status if provided (validate against madmin constants)
	if req.NewStatus != "" {
		if req.NewStatus != string(madmin.AccountEnabled) && req.NewStatus != string(madmin.AccountDisabled) {
			logger.Error().
				Str("accessKey", req.AccessKey).
				Str("invalidStatus", req.NewStatus).
				Msg("Invalid account status provided")
			return nil, fmt.Errorf("invalid account status: %s. Must be 'enabled' or 'disabled'", req.NewStatus)
		}
		updateReq.NewStatus = req.NewStatus
	}

	// Set name if provided
	if req.NewName != "" {
		updateReq.NewName = req.NewName
	}

	// Set description if provided
	if req.NewDescription != "" {
		updateReq.NewDescription = req.NewDescription
	}

	// Set expiration if provided
	if req.NewExpiration != nil && *req.NewExpiration > 0 {
		expTime := time.Unix(*req.NewExpiration, 0)
		updateReq.NewExpiration = &expTime
		logger.Debug().Int64("newExpiration", *req.NewExpiration).Time("newExpirationTime", expTime).Msg("Setting new expiration time")
	}

	// Execute the update
	err := s.minioClient.UpdateServiceAccount(ctx, req.AccessKey, updateReq)
	if err != nil {
		logger.Error().
			Err(err).
			Str("accessKey", req.AccessKey).
			Msg("Failed to update service account")
		return nil, fmt.Errorf("failed to update service account: %w", err)
	}

	logger.Info().
		Str("accessKey", req.AccessKey).
		Msg("Service account updated successfully")

	response := &UpdateServiceAccountResponse{
		AccessKey: req.AccessKey,
		Message:   "Service account updated successfully",
	}

	// Include the new secret key in response if it was rotated
	if req.NewSecretKey != "" {
		response.NewSecretKey = req.NewSecretKey
		logger.Debug().
			Str("accessKey", req.AccessKey).
			Msg("Secret key rotated, including new secret key in response")
	}

	return response, nil
}
