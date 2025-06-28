package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/elct9620/minio-lite-admin/internal/service"
	"github.com/go-chi/chi/v5"
)

// DeleteAccessKeysHandler handles DELETE /api/access-keys/{accessKey} requests
func (s *Service) DeleteAccessKeysHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.logger.With().Str("handler", "DeleteAccessKeysHandler").Logger()

	// Get access key from URL path
	accessKey := chi.URLParam(r, "accessKey")
	if accessKey == "" {
		logger.Warn().Msg("Access key parameter is missing")
		http.Error(w, "Access key is required", http.StatusBadRequest)
		return
	}

	// Validate access key format (basic validation)
	accessKey = strings.TrimSpace(accessKey)
	if accessKey == "" {
		logger.Warn().Msg("Access key parameter is empty")
		http.Error(w, "Access key is required", http.StatusBadRequest)
		return
	}

	logger.Debug().Str("accessKey", accessKey).Msg("Deleting access key")

	// Create service request
	req := service.DeleteServiceAccountRequest{
		AccessKey: accessKey,
	}

	// Execute the service
	response, err := s.deleteServiceAccountService.Execute(ctx, req)
	if err != nil {
		logger.Error().Err(err).Str("accessKey", accessKey).Msg("Failed to delete access key")
		http.Error(w, "Failed to delete access key", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logger.Info().
		Str("accessKey", response.AccessKey).
		Msg("Successfully deleted access key")
}
