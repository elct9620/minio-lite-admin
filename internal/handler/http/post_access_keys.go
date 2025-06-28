package http

import (
	"encoding/json"
	"net/http"

	"github.com/elct9620/minio-lite-admin/internal/service"
	"github.com/rs/zerolog"
)

// PostAccessKeysHandler handles POST /api/access-keys to create a new service account
func (s *Service) PostAccessKeysHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := zerolog.Ctx(ctx)

	// Parse request body
	var req service.CreateServiceAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error().Err(err).Msg("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" {
		logger.Warn().Msg("Service account name is required")
		http.Error(w, "Service account name is required", http.StatusBadRequest)
		return
	}

	// Create service account
	response, err := s.addServiceAccountService.Execute(ctx, req)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create service account")
		http.Error(w, "Failed to create service account", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logger.Info().
		Str("accessKey", response.AccessKey).
		Str("name", response.Name).
		Msg("Successfully created service account")
}
