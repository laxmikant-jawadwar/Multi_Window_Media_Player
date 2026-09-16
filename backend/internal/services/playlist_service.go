package services

import (
	"fmt"
	"media-sequencer/backend/internal/repositories"
)

type PlaylistService struct {
	Repository *repositories.PlaylistRepository
}

func NewPlaylistService(repository *repositories.PlaylistRepository) *PlaylistService {
	return &PlaylistService{
		Repository: repository,
	}
}

func (s *PlaylistService) AddMedia(
	windowID string,
	mediaID string,
	position int,
) (string, error) {
	return s.Repository.AddMedia(windowID, mediaID, position)
}

func (s *PlaylistService) GetByWindow(
	windowID string,
) ([]map[string]interface{}, error) {
	return s.Repository.GetByWindow(windowID)
}

func (s *PlaylistService) Delete(
	windowID string,
	playlistItemID string,
) error {
	return s.Repository.Delete(windowID, playlistItemID)
}

func (s *PlaylistService) UpdatePosition(
	windowID string,
	playlistItemID string,
	newPosition int,
) error {
	if newPosition <= 0 {
		return fmt.Errorf("position must be greater than 0")
	}
	return s.Repository.UpdatePosition(windowID, playlistItemID, newPosition)
}
