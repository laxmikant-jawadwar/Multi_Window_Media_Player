package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"media-sequencer/backend/internal/services"
)

type PlaylistHandler struct {
	Service *services.PlaylistService
}

func NewPlaylistHandler(service *services.PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{
		Service: service,
	}
}

func (h *PlaylistHandler) AddMedia(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 3 {
		http.Error(w, "window id required", http.StatusBadRequest)
		return
	}

	windowID := parts[2]

	var request struct {
		MediaID  string `json:"media_id"`
		Position int    `json:"position"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if request.MediaID == "" || request.Position <= 0 {
		http.Error(w, "media_id and valid position required", http.StatusBadRequest)
		return
	}

	id, err := h.Service.AddMedia(
		windowID,
		request.MediaID,
		request.Position,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":        id,
		"window_id": windowID,
		"media_id":  request.MediaID,
		"position":  request.Position,
	})
}

func (h *PlaylistHandler) GetPlaylist(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 3 {
		http.Error(w, "window id required", http.StatusBadRequest)
		return
	}

	windowID := parts[2]

	playlist, err := h.Service.GetByWindow(windowID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(playlist)
}
