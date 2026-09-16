package handlers

import (
	"encoding/json"
	"net/http"

	"media-sequencer/backend/internal/services"
)

type WindowHandler struct {
	Service *services.WindowService
}

func NewWindowHandler(service *services.WindowService) *WindowHandler {
	return &WindowHandler{
		Service: service,
	}
}

func (h *WindowHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if request.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	id, err := h.Service.Create(request.Name)
	if err != nil {
		http.Error(w, "failed to create window", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"id":   id,
		"name": request.Name,
	})
}

func (h *WindowHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	windows, err := h.Service.GetAll()
	if err != nil {
		http.Error(w, "failed to get windows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(windows)
}
