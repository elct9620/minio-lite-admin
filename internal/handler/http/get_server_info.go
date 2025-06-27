package http

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
)

// ServerInfoResponse represents the MinIO server info response
type ServerInfoResponse struct {
	Mode         string `json:"mode"`
	Region       string `json:"region"`
	DeploymentID string `json:"deploymentId"`
}

// GetServerInfoHandler handles MinIO server info requests
func (s *Service) GetServerInfoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := zerolog.Ctx(ctx)

	logger.Info().Msg("Fetching MinIO server information")

	w.Header().Set("Content-Type", "application/json")

	info, err := s.getServerInfoService.Execute(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get MinIO server info")
		http.Error(w, "Failed to get MinIO server info", http.StatusInternalServerError)
		return
	}

	response := ServerInfoResponse{
		Mode:         info.ServerInfo.Mode,
		Region:       info.ServerInfo.Region,
		DeploymentID: info.ServerInfo.DeploymentID,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error().Err(err).Msg("Failed to encode MinIO server info response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
