package main

import (
	"log"
	"net/http"

	"media-sequencer/backend/internal/config"
	"media-sequencer/backend/internal/database"
	"media-sequencer/backend/internal/handlers"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Printf("Database connected successfully")

	log.Printf("Server starting on port %s", cfg.Port)

	http.HandleFunc("/health",handlers.HealthHandler)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatal(err)
	}
}
