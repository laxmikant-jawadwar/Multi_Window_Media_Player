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
	syncRepository := repositories.NewSyncRepository(db)

	//services
	windowService := services.NewWindowService(windowRepository)
	mediaService := services.NewMediaService(mediaRepository)

	playlistService := services.NewPlaylistService(playlistRepository)
	playbackService := services.NewPlaybackService(
		playbackRepository,
		playlistRepository,
		mediaRepository,
		syncRepository,
	)

	syncService := services.NewSyncService(syncRepository)

	//handlers
	windowHandler := handlers.NewWindowHandler(windowService)
	mediaHandler := handlers.NewMediaHandler(mediaService)

	playlistHandler := handlers.NewPlaylistHandler(playlistService)

	playbackHandler := handlers.NewPlaybackHandler(playbackService)

	syncHandler := handlers.NewSyncHandler(syncService)

	http.HandleFunc("/", handlers.WelcomeHandler)
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

	//sync route
	http.HandleFunc("/sync", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			syncHandler.Start(w, r)
			return
		}

		if r.Method == http.MethodGet {
			syncHandler.GetStatus(w, r)
			return
		}

		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	// Window specific routes
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

			if r.Method == http.MethodDelete {
				playlistHandler.DeleteItem(w, r)
				return
			}

			if r.Method == http.MethodPut {
				playlistHandler.UpdatePosition(w, r)
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

		// Playback current
		if strings.HasSuffix(path, "/playback/current") {

			if r.Method == http.MethodGet {
				playbackHandler.GetCurrentMedia(w, r)
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

	handler := enableCORS(http.DefaultServeMux)

	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatal(err)
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
