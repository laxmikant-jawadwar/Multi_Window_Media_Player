package services

import "media-sequencer/backend/internal/repositories"

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
