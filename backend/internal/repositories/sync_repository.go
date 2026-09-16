package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type SyncRepository struct {
	DB *sql.DB
}

func NewSyncRepository(db *sql.DB) *SyncRepository {
	return &SyncRepository{DB: db}
}

func (r *SyncRepository) Start(mediaID string, durationSeconds int) (string, error) {
	id := uuid.New()

	_, err := r.DB.Exec(`
		INSERT INTO sync_sessions
			(id, media_id, started_at, duration_seconds)
		VALUES ($1, $2, $3, $4)
	`, id, mediaID, time.Now().UTC(), durationSeconds)

	if err != nil {
		return "", fmt.Errorf("failed to start sync session: %w", err)
	}

	return id.String(), nil
}

func (r *SyncRepository) GetLatest() (string, string, time.Time, int, error) {
	var (
		sessionID       string
		mediaID         string
		startedAt       time.Time
		durationSeconds int
	)

	err := r.DB.QueryRow(`
		SELECT id, media_id, started_at, duration_seconds
		FROM sync_sessions
		ORDER BY started_at DESC
		LIMIT 1
	`).Scan(
		&sessionID,
		&mediaID,
		&startedAt,
		&durationSeconds,
	)

	if err != nil {
		return "", "", time.Time{}, 0, err
	}

	return sessionID, mediaID, startedAt, durationSeconds, nil
}
