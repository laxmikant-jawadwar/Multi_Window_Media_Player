package repositories

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type WindowRepository struct {
	DB *sql.DB
}

func NewWindowRepository(db *sql.DB) *WindowRepository {
	return &WindowRepository{DB: db}
}

func (r *WindowRepository) Create(name string) (string, error) {
	id := uuid.New()

	_, err := r.DB.Exec(
		`INSERT INTO windows (id, name) VALUES ($1, $2)`,
		id,
		name,
	)

	if err != nil {
		return "", fmt.Errorf("failed to create window: %w", err)
	}

	return id.String(), nil
}

func (r *WindowRepository) GetAll() ([]map[string]interface{}, error) {
	rows, err := r.DB.Query(
		`SELECT id, name, created_at, updated_at FROM windows ORDER BY created_at`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var windows []map[string]interface{}

	for rows.Next() {
		var id, name string
		var createdAt, updatedAt interface{}

		if err := rows.Scan(&id, &name, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		windows = append(windows, map[string]interface{}{
			"id":         id,
			"name":       name,
			"created_at": createdAt,
			"updated_at": updatedAt,
		})
	}

	return windows, nil
}
