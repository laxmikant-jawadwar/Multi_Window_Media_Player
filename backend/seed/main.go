package main

import (
	"log"

	"media-sequencer/backend/internal/config"
	"media-sequencer/backend/internal/database"
)

type WindowSeed struct {
	ID   string
	Name string
}

type MediaSeed struct {
	ID              string
	Title           string
	MediaType       string
	URL             string
	DurationSeconds int
}

type PlaylistItemSeed struct {
	ID       string
	WindowID string
	MediaID  string
	Position int
}

func main() {
	cfg := config.Load()

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to PostgreSQL database for seeding...")

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	//Seed Windows, taken 3 windows
	windows := []WindowSeed{
		{ID: "a0000000-0000-0000-0000-000000000001", Name: "Window 1"},
		{ID: "a0000000-0000-0000-0000-000000000002", Name: "Window 2"},
		{ID: "a0000000-0000-0000-0000-000000000003", Name: "Window 3"},
	}

	for _, w := range windows {
		_, err := tx.Exec(`
			INSERT INTO windows (id, name)
			VALUES ($1, $2)
			ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name
		`, w.ID, w.Name)
		if err != nil {
			log.Fatalf("Failed to seed window %s: %v", w.Name, err)
		}
	}
	log.Printf("Seeded %d windows", len(windows))

	//Seed Media Records (5 Samsung Demo Media)
	mediaList := []MediaSeed{
		{
			ID:              "b0000000-0000-0000-0000-000000000001",
			Title:           "Galaxy S26 Product Video",
			MediaType:       "video",
			URL:             "https://example.com/media/galaxy-s26-video.mp4",
			DurationSeconds: 30,
		},
		{
			ID:              "b0000000-0000-0000-0000-000000000002",
			Title:           "Galaxy S26 Product Image",
			MediaType:       "image",
			URL:             "https://example.com/media/galaxy-s26-image.jpg",
			DurationSeconds: 10,
		},
		{
			ID:              "b0000000-0000-0000-0000-000000000003",
			Title:           "Galaxy Buds Product Video",
			MediaType:       "video",
			URL:             "https://example.com/media/galaxy-buds-video.mp4",
			DurationSeconds: 25,
		},
		{
			ID:              "b0000000-0000-0000-0000-000000000004",
			Title:           "Galaxy Watch Product Image",
			MediaType:       "image",
			URL:             "https://example.com/media/galaxy-watch-image.jpg",
			DurationSeconds: 10,
		},
		{
			ID:              "b0000000-0000-0000-0000-000000000005",
			Title:           "Samsung Ecosystem Video",
			MediaType:       "video",
			URL:             "https://example.com/media/samsung-ecosystem-video.mp4",
			DurationSeconds: 35,
		},
	}

	for _, m := range mediaList {
		_, err := tx.Exec(`
			INSERT INTO media (id, title, media_type, url, duration_seconds)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (id) DO UPDATE SET
				title = EXCLUDED.title,
				media_type = EXCLUDED.media_type,
				url = EXCLUDED.url,
				duration_seconds = EXCLUDED.duration_seconds
		`, m.ID, m.Title, m.MediaType, m.URL, m.DurationSeconds)
		if err != nil {
			log.Fatalf("Failed to seed media %s: %v", m.Title, err)
		}
	}
	log.Printf("Seeded %d media records", len(mediaList))

	// Seed Playlist Assignments
	w1 := "a0000000-0000-0000-0000-000000000001"
	w2 := "a0000000-0000-0000-0000-000000000002"
	w3 := "a0000000-0000-0000-0000-000000000003"

	m1 := "b0000000-0000-0000-0000-000000000001"
	m2 := "b0000000-0000-0000-0000-000000000002"
	m3 := "b0000000-0000-0000-0000-000000000003"
	m4 := "b0000000-0000-0000-0000-000000000004"
	m5 := "b0000000-0000-0000-0000-000000000005"

	playlistItems := []PlaylistItemSeed{
		// Window 1
		{ID: "c0000000-0000-0000-0001-000000000001", WindowID: w1, MediaID: m1, Position: 1},
		{ID: "c0000000-0000-0000-0001-000000000002", WindowID: w1, MediaID: m2, Position: 2},
		{ID: "c0000000-0000-0000-0001-000000000003", WindowID: w1, MediaID: m5, Position: 3},

		// Window 2
		{ID: "c0000000-0000-0000-0002-000000000001", WindowID: w2, MediaID: m3, Position: 1},
		{ID: "c0000000-0000-0000-0002-000000000002", WindowID: w2, MediaID: m4, Position: 2},
		{ID: "c0000000-0000-0000-0002-000000000003", WindowID: w2, MediaID: m1, Position: 3},

		// Window 3
		{ID: "c0000000-0000-0000-0003-000000000001", WindowID: w3, MediaID: m4, Position: 1},
		{ID: "c0000000-0000-0000-0003-000000000002", WindowID: w3, MediaID: m5, Position: 2},
		{ID: "c0000000-0000-0000-0003-000000000003", WindowID: w3, MediaID: m3, Position: 3},
	}

	for _, pi := range playlistItems {
		_, err := tx.Exec(`
			INSERT INTO playlist_items (id, window_id, media_id, position)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (window_id, position) DO UPDATE SET
				id = EXCLUDED.id,
				media_id = EXCLUDED.media_id
		`, pi.ID, pi.WindowID, pi.MediaID, pi.Position)
		if err != nil {
			log.Fatalf("Failed to seed playlist item for window %s position %d: %v", pi.WindowID, pi.Position, err)
		}
	}
	log.Printf("Seeded %d playlist assignments", len(playlistItems))

	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit seed transaction: %v", err)
	}

	log.Println("Database seeding completed successfully!")
}
