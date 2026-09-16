package services

import (
	"fmt"
	"time"

	"media-sequencer/backend/internal/repositories"
)

const FiveHourCycle = 5 * time.Hour

type PlaybackService struct {
	PlaybackRepository *repositories.PlaybackRepository
	PlaylistRepository *repositories.PlaylistRepository
}

func NewPlaybackService(
	playbackRepository *repositories.PlaybackRepository,
	playlistRepository *repositories.PlaylistRepository,
) *PlaybackService {
	return &PlaybackService{
		PlaybackRepository: playbackRepository,
		PlaylistRepository: playlistRepository,
	}
}

func (s *PlaybackService) Start(windowID string) error {
	return s.PlaybackRepository.Start(windowID)
}

func (s *PlaybackService) Stop(windowID string) error {
	return s.PlaybackRepository.Stop(windowID)
}

func (s *PlaybackService) GetState(windowID string) (map[string]interface{}, error) {

	itemID, startedAt, status, err :=
		s.PlaybackRepository.GetState(windowID)

	if err != nil {
		return nil, err
	}

	if startedAt == nil {
		return nil, fmt.Errorf("playback has not been started")
	}

	elapsed := time.Since(*startedAt)

	if elapsed >= FiveHourCycle && status == "RUNNING" {
		_ = s.PlaybackRepository.Stop(windowID)
		status = "STOPPED"
	}

	return map[string]interface{}{
		"window_id":        windowID,
		"status":           status,
		"current_item_id":  itemID,
		"cycle_started_at": *startedAt,
		"elapsed_seconds":  int(elapsed.Seconds()),
		"remaining_seconds": max(
			int(FiveHourCycle.Seconds()-elapsed.Seconds()),
			0,
		),
	}, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (s *PlaybackService) GetCurrentMedia(windowID string) (map[string]interface{}, error) {
	_, startedAt, status, err :=
		s.PlaybackRepository.GetState(windowID)

	if err != nil {
		return nil, err
	}

	if startedAt == nil {
		return nil, fmt.Errorf("playback has not been started")
	}

	elapsed := time.Since(*startedAt)

	// Stop after 5 hours
	if elapsed >= FiveHourCycle {

		if status == "RUNNING" {
			_ = s.PlaybackRepository.Stop(windowID)
		}

		return map[string]interface{}{
			"window_id": windowID,
			"status":    "STOPPED",
			"message":   "5-hour playback cycle completed",
		}, nil
	}

	playlist, err :=
		s.PlaylistRepository.GetPlaylistWithDurations(windowID)

	if err != nil {
		return nil, err
	}

	if len(playlist) == 0 {
		return nil, fmt.Errorf("playlist is empty")
	}

	var totalDuration int

	for _, item := range playlist {
		totalDuration += item.DurationSeconds
	}

	if totalDuration <= 0 {
		return nil, fmt.Errorf("invalid playlist duration")
	}

	positionInCycle :=
		int(elapsed.Seconds()) % totalDuration

	accumulated := 0

	for _, item := range playlist {

		accumulated += item.DurationSeconds

		if positionInCycle < accumulated {

			_ = s.PlaybackRepository.UpdateCurrentItem(
				windowID,
				item.ID,
			)

			return map[string]interface{}{
				"window_id":        windowID,
				"status":           "RUNNING",
				"media_id":         item.MediaID,
				"title":            item.Title,
				"media_type":       item.MediaType,
				"url":              item.URL,
				"position":         item.Position,
				"duration_seconds": item.DurationSeconds,
				"elapsed_seconds":  int(elapsed.Seconds()),
				"remaining_seconds": int(
					FiveHourCycle.Seconds() - elapsed.Seconds(),
				),
			}, nil
		}
	}

	return nil, fmt.Errorf("unable to determine current media")
}
