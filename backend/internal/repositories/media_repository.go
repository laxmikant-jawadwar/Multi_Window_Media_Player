package repositories

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type MediaRepository struct {
	DB *sql.DB
}

func NewMediaRepository(db *sql.DB) *MediaRepository {
	return &MediaRepository{DB: db}
}

func (r *MediaRepository) Create(
	title string,
	mediaType string,
	url string,
	durationSeconds int,
) (string, error) {

	id := uuid.New()

	_, err := r.DB.Exec(
		`INSERT INTO media
		(id, title, media_type, url, duration_seconds)
		VALUES ($1, $2, $3, $4, $5)`,
		id,
		title,
		mediaType,
		url,
		durationSeconds,
	)

	if err != nil {
		return "", fmt.Errorf("failed to create media: %w", err)
	}

	return id.String(), nil
}

func (r *MediaRepository) GetAll() ([]map[string]interface{}, error) {

	rows, err := r.DB.Query(`
		SELECT id, title, media_type, url, duration_seconds,
		       created_at, updated_at
		FROM media
		ORDER BY created_at
	`)

	if err != nil {
		return nil, fmt.Errorf("failed to get media: %w", err)
	}

	defer rows.Close()

	var mediaList []map[string]interface{}

	for rows.Next() {

		var (
			id              string
			title           string
			mediaType       string
			url             string
			durationSeconds int
			createdAt       interface{}
			updatedAt       interface{}
		)

		if err := rows.Scan(
			&id,
			&title,
			&mediaType,
			&url,
			&durationSeconds,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan media: %w", err)
		}

		mediaList = append(mediaList, map[string]interface{}{
			"id":               id,
			"title":            title,
			"media_type":       mediaType,
			"url":              url,
			"duration_seconds": durationSeconds,
			"created_at":       createdAt,
			"updated_at":       updatedAt,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed while reading media: %w", err)
	}

	return mediaList, nil
}

func (r *MediaRepository) GetByID(mediaID string) (map[string]interface{}, error) {
	var (
		id              string
		title           string
		mediaType       string
		url             string
		durationSeconds int
		createdAt       interface{}
		updatedAt       interface{}
	)

	err := r.DB.QueryRow(`
		SELECT id, title, media_type, url, duration_seconds,
		       created_at, updated_at
		FROM media
		WHERE id = $1
	`, mediaID).Scan(
		&id,
		&title,
		&mediaType,
		&url,
		&durationSeconds,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get media: %w", err)
	}

	return map[string]interface{}{
		"id":               id,
		"title":            title,
		"media_type":       mediaType,
		"url":              url,
		"duration_seconds": durationSeconds,
		"created_at":       createdAt,
		"updated_at":       updatedAt,
	}, nil
}
