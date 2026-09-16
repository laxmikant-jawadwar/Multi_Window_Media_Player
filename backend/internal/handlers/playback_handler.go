package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"media-sequencer/backend/internal/services"
)

type PlaybackHandler struct {
	Service *services.PlaybackService
}

func NewPlaybackHandler(service *services.PlaybackService) *PlaybackHandler {
	return &PlaybackHandler{
		Service: service,
	}
}

func (h *PlaybackHandler) Start(w http.ResponseWriter, r *http.Request) {

	windowID := strings.TrimPrefix(r.URL.Path, "/windows/")
	windowID = strings.TrimSuffix(windowID, "/playback/start")

	if err := h.Service.Start(windowID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"window_id": windowID,
		"status":    "RUNNING",
	})
}

func (h *PlaybackHandler) GetState(w http.ResponseWriter, r *http.Request) {

	windowID := strings.TrimPrefix(r.URL.Path, "/windows/")
	windowID = strings.TrimSuffix(windowID, "/playback")

	state, err := h.Service.GetState(windowID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(state)
}

func (h *PlaybackHandler) GetCurrentMedia(w http.ResponseWriter, r *http.Request) {

	windowID := strings.TrimPrefix(r.URL.Path, "/windows/")
	windowID = strings.TrimSuffix(windowID, "/playback/current")

	media, err := h.Service.GetCurrentMedia(windowID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(media)
}
