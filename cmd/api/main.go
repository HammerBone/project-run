package main

import (
	"log"

	"github.com/HammerBone/project-run/internal/config"
	"github.com/HammerBone/project-run/internal/database"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := config.LoadConfig("config/development.yaml")
	if err != nil {
		log.Fatal(err)
	}

	_, err = database.NewDBConnection(cfg.DB.DBAddr)
	if err != nil {
		log.Fatal(err)
	}

	app := &application{
		config: cfg,
		dbConfig: cfg.DB,
	}

	mux := chi.NewRouter()
	app.run(mux)
}
