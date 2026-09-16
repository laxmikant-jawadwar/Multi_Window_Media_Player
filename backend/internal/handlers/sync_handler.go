package handlers

import (
	"encoding/json"
	"net/http"

	"media-sequencer/backend/internal/services"
)

type SyncHandler struct {
	Service *services.SyncService
}

func NewSyncHandler(service *services.SyncService) *SyncHandler {
	return &SyncHandler{
		Service: service,
	}
}

func (h *SyncHandler) Start(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		MediaID         string `json:"media_id"`
		DurationSeconds int    `json:"duration_seconds"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	sessionID, err := h.Service.Start(
		request.MediaID,
		request.DurationSeconds,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"session_id":       sessionID,
		"media_id":         request.MediaID,
		"duration_seconds": request.DurationSeconds,
		"status":           "ACTIVE",
	})
}

func (h *SyncHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status, err := h.Service.GetStatus()

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(status)
}
