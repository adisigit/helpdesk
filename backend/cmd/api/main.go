package main

import (
	"helpdesk/backend/internal/config"
	"helpdesk/backend/internal/database"
	"helpdesk/backend/internal/server"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	if err := database.RunMigrations(cfg.DatabaseURL, cfg.MigrationsPath); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}

	if cfg.SeedOnStart {
		if err := database.Seed(db, cfg.SeedPassword); err != nil {
			log.Fatalf("seed: %v", err)
		}
	}

	srv := server.New(cfg, db)
	log.Printf("listening on :%s (docs: /docs)", cfg.Port)
	log.Fatal(srv.ListenAndServe())
}
