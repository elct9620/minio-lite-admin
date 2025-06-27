package http

import (
	"encoding/json"
	"net/http"

	"github.com/elct9620/minio-lite-admin/internal/service"
	"github.com/rs/zerolog"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// ServerInfoResponse represents the server info response
type ServerInfoResponse struct {
	Version string `json:"version"`
	Name    string `json:"name"`
}

// MinIOServerInfoResponse represents the MinIO server info response
type MinIOServerInfoResponse struct {
	Mode         string `json:"mode"`
	Region       string `json:"region"`
	DeploymentID string `json:"deploymentId"`
}

// HealthHandler handles health check requests
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := HealthResponse{
		Status:  "ok",
		Service: "minio-lite-admin",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ServerInfoHandler handles server info requests
func ServerInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := ServerInfoResponse{
		Version: "0.1.0",
		Name:    "MinIO Lite Admin",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// MinIOServerInfoHandler handles MinIO server info requests
func MinIOServerInfoHandler(getServerInfoService *service.GetServerInfoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := zerolog.Ctx(ctx)

		logger.Info().Msg("Fetching MinIO server information")

		w.Header().Set("Content-Type", "application/json")

		info, err := getServerInfoService.Execute(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to get MinIO server info")
			http.Error(w, "Failed to get MinIO server info", http.StatusInternalServerError)
			return
		}

		response := MinIOServerInfoResponse{
			Mode:         info.Mode,
			Region:       info.Region,
			DeploymentID: info.DeploymentID,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Error().Err(err).Msg("Failed to encode MinIO server info response")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
