package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PlaybackRepository struct {
	DB *sql.DB
}

func NewPlaybackRepository(db *sql.DB) *PlaybackRepository {
	return &PlaybackRepository{DB: db}
}

func (r *PlaybackRepository) Start(windowID string) error {
	_, err := r.DB.Exec(`
		INSERT INTO playback_states
			(window_id, cycle_started_at, status, updated_at)
		VALUES ($1, $2, 'RUNNING', $2)
		ON CONFLICT (window_id)
		DO UPDATE SET
			cycle_started_at = EXCLUDED.cycle_started_at,
			status = 'RUNNING',
			updated_at = EXCLUDED.updated_at
	`, windowID, time.Now().UTC())

	return err
}

func (r *PlaybackRepository) Stop(windowID string) error {
	_, err := r.DB.Exec(`
		UPDATE playback_states
		SET status = 'STOPPED',
		    updated_at = $2
		WHERE window_id = $1
	`, windowID, time.Now().UTC())

	return err
}

func (r *PlaybackRepository) GetState(windowID string) (
	string, *time.Time, string, error,
) {
	var (
		startedAt *time.Time
		status    string
		itemID    sql.NullString
	)

	err := r.DB.QueryRow(`
		SELECT cycle_started_at, status, current_playlist_item_id
		FROM playback_states
		WHERE window_id = $1
	`, windowID).Scan(&startedAt, &status, &itemID)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil, "STOPPED", nil
		}
		return "", nil, "", fmt.Errorf("failed to get playback state: %w", err)
	}

	return itemID.String, startedAt, status, nil
}

func (r *PlaybackRepository) UpdateCurrentItem(
	windowID string,
	itemID string,
) error {

	_, err := r.DB.Exec(`
		UPDATE playback_states
		SET current_playlist_item_id = $2,
		    updated_at = $3
		WHERE window_id = $1
	`, windowID, uuid.MustParse(itemID), time.Now().UTC())

	return err
}
