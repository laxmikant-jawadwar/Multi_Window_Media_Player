package repositories

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type PlaylistRepository struct {
	DB *sql.DB
}

func NewPlaylistRepository(db *sql.DB) *PlaylistRepository {
	return &PlaylistRepository{DB: db}
}

func (r *PlaylistRepository) AddMedia(
	windowID string,
	mediaID string,
	position int,
) (string, error) {

	id := uuid.New()

	_, err := r.DB.Exec(`
		INSERT INTO playlist_items
		(id, window_id, media_id, position)
		VALUES ($1, $2, $3, $4)
	`, id, windowID, mediaID, position)

	if err != nil {
		return "", fmt.Errorf("failed to add media to playlist: %w", err)
	}

	return id.String(), nil
}

func (r *PlaylistRepository) GetByWindow(
	windowID string,
) ([]map[string]interface{}, error) {

	rows, err := r.DB.Query(`
		SELECT
			pi.id,
			pi.position,
			m.id,
			m.title,
			m.media_type,
			m.url,
			m.duration_seconds
		FROM playlist_items pi
		JOIN media m ON pi.media_id = m.id
		WHERE pi.window_id = $1
		ORDER BY pi.position
	`, windowID)

	if err != nil {
		return nil, fmt.Errorf("failed to get playlist: %w", err)
	}

	defer rows.Close()

	var playlist []map[string]interface{}

	for rows.Next() {
		var (
			itemID          string
			mediaID         string
			title           string
			mediaType       string
			url             string
			position        int
			durationSeconds int
		)

		if err := rows.Scan(
			&itemID,
			&position,
			&mediaID,
			&title,
			&mediaType,
			&url,
			&durationSeconds,
		); err != nil {
			return nil, err
		}

		playlist = append(playlist, map[string]interface{}{
			"id":               itemID,
			"media_id":         mediaID,
			"title":            title,
			"media_type":       mediaType,
			"url":              url,
			"position":         position,
			"duration_seconds": durationSeconds,
		})
	}

	return playlist, nil
}

func (r *PlaylistRepository) GetPlaylistWithDurations(
	windowID string,
) ([]struct {
	ID              string
	MediaID         string
	Title           string
	MediaType       string
	URL             string
	Position        int
	DurationSeconds int
}, error) {

	rows, err := r.DB.Query(`
		SELECT
			pi.id,
			m.id,
			m.title,
			m.media_type,
			m.url,
			pi.position,
			m.duration_seconds
		FROM playlist_items pi
		JOIN media m ON pi.media_id = m.id
		WHERE pi.window_id = $1
		ORDER BY pi.position
	`, windowID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var playlist []struct {
		ID              string
		MediaID         string
		Title           string
		MediaType       string
		URL             string
		Position        int
		DurationSeconds int
	}

	for rows.Next() {

		var item struct {
			ID              string
			MediaID         string
			Title           string
			MediaType       string
			URL             string
			Position        int
			DurationSeconds int
		}

		err := rows.Scan(
			&item.ID,
			&item.MediaID,
			&item.Title,
			&item.MediaType,
			&item.URL,
			&item.Position,
			&item.DurationSeconds,
		)

		if err != nil {
			return nil, err
		}

		playlist = append(playlist, item)
	}

	return playlist, nil
}
