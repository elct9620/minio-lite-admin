package service

import (
	"context"
	"fmt"

	"github.com/minio/madmin-go/v4"
	"github.com/rs/zerolog"
)

type GetServerInfoService struct {
	minioClient *madmin.AdminClient
}

func NewGetServerInfoService(minioClient *madmin.AdminClient) *GetServerInfoService {
	return &GetServerInfoService{
		minioClient: minioClient,
	}
}

func (s *GetServerInfoService) Execute(ctx context.Context) (madmin.InfoMessage, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("Fetching MinIO server info")

	info, err := s.minioClient.ServerInfo(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to fetch MinIO server info")
		return madmin.InfoMessage{}, fmt.Errorf("failed to get server info: %w", err)
	}

	logger.Debug().Msg("Successfully fetched MinIO server info")
	return info, nil
}
