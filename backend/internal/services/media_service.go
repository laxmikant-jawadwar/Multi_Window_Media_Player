package services

import "media-sequencer/backend/internal/repositories"

type MediaService struct {
	Repository *repositories.MediaRepository
}

func NewMediaService(repository *repositories.MediaRepository) *MediaService {
	return &MediaService{
		Repository: repository,
	}
}

func (s *MediaService) Create(
	title string,
	mediaType string,
	url string,
	durationSeconds int,
) (string, error) {

	return s.Repository.Create(
		title,
		mediaType,
		url,
		durationSeconds,
	)
}

func (s *MediaService) GetAll() ([]map[string]interface{}, error) {
	return s.Repository.GetAll()
}
