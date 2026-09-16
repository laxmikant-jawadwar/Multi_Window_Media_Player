package main

import (
	"log"
	"media-sequencer/backend/internal/config"
	"media-sequencer/backend/internal/database"
	"media-sequencer/backend/internal/handlers"
	"media-sequencer/backend/internal/repositories"
	"media-sequencer/backend/internal/services"
	"net/http"
	"strings"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Printf("Database connected successfully")

	//repos
	windowRepository := repositories.NewWindowRepository(db)
	mediaRepository := repositories.NewMediaRepository(db)
	playlistRepository := repositories.NewPlaylistRepository(db)
	playbackRepository := repositories.NewPlaybackRepository(db)

	//services
	windowService := services.NewWindowService(windowRepository)
	mediaService := services.NewMediaService(mediaRepository)

	playlistService := services.NewPlaylistService(playlistRepository)

	playbackService := services.NewPlaybackService(
		playbackRepository,
		playlistRepository,
	)

	//handlers
	windowHandler := handlers.NewWindowHandler(windowService)
	mediaHandler := handlers.NewMediaHandler(mediaService)

	playlistHandler := handlers.NewPlaylistHandler(playlistService)

	playbackHandler := handlers.NewPlaybackHandler(playbackService)

	http.HandleFunc("/health", handlers.HealthHandler)

	http.HandleFunc("/windows", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodPost {
			windowHandler.Create(w, r)
			return
		}

		if r.Method == http.MethodGet {
			windowHandler.GetAll(w, r)
			return
		}

		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	// Media
	http.HandleFunc("/media", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodPost {
			mediaHandler.Create(w, r)
			return
		}

		if r.Method == http.MethodGet {
			mediaHandler.GetAll(w, r)
			return
		}

		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	// Window-specific routes
	http.HandleFunc("/windows/", func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.Path

		// Playlist routes
		if strings.Contains(path, "/playlist") {

			if r.Method == http.MethodPost {
				playlistHandler.AddMedia(w, r)
				return
			}

			if r.Method == http.MethodGet {
				playlistHandler.GetPlaylist(w, r)
				return
			}
		}

		// Playback start
		if strings.HasSuffix(path, "/playback/start") {

			if r.Method == http.MethodPost {
				playbackHandler.Start(w, r)
				return
			}
		}

		// Playback state
		if strings.HasSuffix(path, "/playback") {

			if r.Method == http.MethodGet {
				playbackHandler.GetState(w, r)
				return
			}
		}

		http.Error(w, "not found", http.StatusNotFound)
	})

	log.Printf("Server starting on port %s", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatal(err)
	}
}