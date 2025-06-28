package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/elct9620/minio-lite-admin/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

func (svc *Service) PutAccessKeysHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := zerolog.Ctx(ctx)

	// Get access key from URL parameter
	accessKey := chi.URLParam(r, "accessKey")
	if strings.TrimSpace(accessKey) == "" {
		logger.Error().Msg("Access key parameter is required")
		http.Error(w, "Access key is required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var updateReq struct {
		NewPolicy      string `json:"newPolicy,omitempty"`
		NewSecretKey   string `json:"newSecretKey,omitempty"`
		NewStatus      string `json:"newStatus,omitempty"`
		NewName        string `json:"newName,omitempty"`
		NewDescription string `json:"newDescription,omitempty"`
		NewExpiration  *int64 `json:"newExpiration,omitempty"` // Unix timestamp in seconds
	}

	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		logger.Error().
			Err(err).
			Str("accessKey", accessKey).
			Msg("Failed to decode update request")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create service request
	serviceReq := service.UpdateServiceAccountRequest{
		AccessKey:      strings.TrimSpace(accessKey),
		NewPolicy:      strings.TrimSpace(updateReq.NewPolicy),
		NewSecretKey:   strings.TrimSpace(updateReq.NewSecretKey),
		NewStatus:      strings.TrimSpace(updateReq.NewStatus),
		NewName:        strings.TrimSpace(updateReq.NewName),
		NewDescription: strings.TrimSpace(updateReq.NewDescription),
		NewExpiration:  updateReq.NewExpiration,
	}

	// Execute service request
	response, err := svc.updateServiceAccountService.Execute(ctx, serviceReq)
	if err != nil {
		logger.Error().
			Err(err).
			Str("accessKey", accessKey).
			Msg("Failed to update service account")
		http.Error(w, "Failed to update access key", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error().
			Err(err).
			Str("accessKey", accessKey).
			Msg("Failed to encode update response")
		return
	}

	logger.Info().
		Str("accessKey", accessKey).
		Msg("Service account updated successfully")
}
