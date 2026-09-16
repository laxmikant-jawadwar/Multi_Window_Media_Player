package models

import "time"

type Media struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	MediaType       string    `json:"media_type"`
	URL             string    `json:"url"`
	DurationSeconds int       `json:"duration_seconds"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
