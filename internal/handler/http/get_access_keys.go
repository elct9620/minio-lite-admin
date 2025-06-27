package http

import (
	"encoding/json"
	"net/http"

	"github.com/elct9620/minio-lite-admin/internal/service"
)

// GetAccessKeysHandler handles GET /api/access-keys requests
func (s *Service) GetAccessKeysHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.logger.With().Str("handler", "GetAccessKeysHandler").Logger()

	// Parse query parameters
	query := r.URL.Query()
	filterType := query.Get("type")
	filterUser := query.Get("user")

	// Default to "all" if no type specified
	if filterType == "" {
		filterType = "all"
	}

	// Validate filter type
	validTypes := map[string]bool{
		"all":             true,
		"users":           true,
		"serviceAccounts": true,
		"sts":             true,
	}

	if !validTypes[filterType] {
		logger.Warn().Str("type", filterType).Msg("Invalid access key type filter")
		http.Error(w, "Invalid type parameter. Valid values: all, users, serviceAccounts, sts", http.StatusBadRequest)
		return
	}

	logger.Debug().Str("type", filterType).Str("user", filterUser).Msg("Getting access keys")

	// Create options for the service
	opts := service.ListAccessKeysOptions{
		Type: filterType,
		User: filterUser,
	}

	// Execute the service
	result, err := s.listAccessKeysService.Execute(ctx, opts)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get access keys")
		http.Error(w, "Failed to get access keys", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Encode and send response
	if err := json.NewEncoder(w).Encode(result); err != nil {
		logger.Error().Err(err).Msg("Failed to encode access keys response")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logger.Debug().Int("total", result.Total).Msg("Successfully returned access keys")
}
