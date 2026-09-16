package services

import "media-sequencer/backend/internal/repositories"

type WindowService struct {
	Repository *repositories.WindowRepository
}

func NewWindowService(repository *repositories.WindowRepository) *WindowService {
	return &WindowService{
		Repository: repository,
	}
}

func (s *WindowService) Create(name string) (string, error) {
	return s.Repository.Create(name)
}

func (s *WindowService) GetAll() ([]map[string]interface{}, error) {
	return s.Repository.GetAll()
}
