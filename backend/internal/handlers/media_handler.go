package handlers

import (
	"encoding/json"
	"net/http"

	"media-sequencer/backend/internal/services"
)

type MediaHandler struct {
	Service *services.MediaService
}

func NewMediaHandler(service *services.MediaService) *MediaHandler {
	return &MediaHandler{
		Service: service,
	}
}

func (h *MediaHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Title           string `json:"title"`
		MediaType       string `json:"media_type"`
		URL             string `json:"url"`
		DurationSeconds int    `json:"duration_seconds"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if request.Title == "" || request.MediaType == "" || request.URL == "" {
		http.Error(w, "title, media_type and url are required", http.StatusBadRequest)
		return
	}

	if request.DurationSeconds <= 0 {
		http.Error(w, "duration_seconds must be greater than 0", http.StatusBadRequest)
		return
	}

	id, err := h.Service.Create(
		request.Title,
		request.MediaType,
		request.URL,
		request.DurationSeconds,
	)

	if err != nil {
		http.Error(w, "failed to create media", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":               id,
		"title":            request.Title,
		"media_type":       request.MediaType,
		"url":              request.URL,
		"duration_seconds": request.DurationSeconds,
	})
}

func (h *MediaHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	media, err := h.Service.GetAll()
	if err != nil {
		http.Error(w, "failed to get media", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(media)
}
