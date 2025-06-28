package service

import (
	"context"
	"fmt"

	"github.com/minio/madmin-go/v4"
	"github.com/rs/zerolog"
)

type DeleteServiceAccountService struct {
	minioClient *madmin.AdminClient
}

// DeleteServiceAccountRequest represents the request to delete a service account
type DeleteServiceAccountRequest struct {
	AccessKey string `json:"accessKey"`
}

// DeleteServiceAccountResponse represents the response from deleting a service account
type DeleteServiceAccountResponse struct {
	AccessKey string `json:"accessKey"`
	Message   string `json:"message"`
}

func NewDeleteServiceAccountService(minioClient *madmin.AdminClient) *DeleteServiceAccountService {
	return &DeleteServiceAccountService{
		minioClient: minioClient,
	}
}

func (s *DeleteServiceAccountService) Execute(ctx context.Context, req DeleteServiceAccountRequest) (*DeleteServiceAccountResponse, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().
		Str("accessKey", req.AccessKey).
		Msg("Deleting service account")

	// Delete the service account
	err := s.minioClient.DeleteServiceAccount(ctx, req.AccessKey)
	if err != nil {
		logger.Error().Err(err).Str("accessKey", req.AccessKey).Msg("Failed to delete service account")
		return nil, fmt.Errorf("failed to delete service account: %w", err)
	}

	logger.Info().
		Str("accessKey", req.AccessKey).
		Msg("Successfully deleted service account")

	// Prepare response
	response := &DeleteServiceAccountResponse{
		AccessKey: req.AccessKey,
		Message:   "Service account deleted successfully",
	}

	return response, nil
}