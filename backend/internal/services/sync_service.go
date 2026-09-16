package services

import (
	"fmt"
	"time"

	"media-sequencer/backend/internal/repositories"
)

type SyncService struct {
	Repository *repositories.SyncRepository
}

func NewSyncService(repository *repositories.SyncRepository) *SyncService {
	return &SyncService{
		Repository: repository,
	}
}

func (s *SyncService) Start(mediaID string, durationSeconds int) (string, error) {
	if mediaID == "" {
		return "", fmt.Errorf("media_id is required")
	}

	if durationSeconds <= 0 {
		return "", fmt.Errorf("duration must be greater than 0")
	}

	return s.Repository.Start(mediaID, durationSeconds)
}

func (s *SyncService) GetStatus() (map[string]interface{}, error) {
	sessionID, mediaID, startedAt, durationSeconds, err :=
		s.Repository.GetLatest()

	if err != nil {
		return nil, fmt.Errorf("no sync session found")
	}

	elapsed := time.Since(startedAt)

	if elapsed >= time.Duration(durationSeconds)*time.Second {
		return map[string]interface{}{
			"status":     "EXPIRED",
			"session_id": sessionID,
			"media_id":   mediaID,
		}, nil
	}

	remaining := durationSeconds - int(elapsed.Seconds())

	return map[string]interface{}{
		"status":            "ACTIVE",
		"session_id":        sessionID,
		"media_id":          mediaID,
		"started_at":        startedAt,
		"duration_seconds":  durationSeconds,
		"elapsed_seconds":   int(elapsed.Seconds()),
		"remaining_seconds": remaining,
	}, nil
}
