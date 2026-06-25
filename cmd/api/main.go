package main

import (
	"log"

	"github.com/HammerBone/project-run/internal/config"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := config.LoadConfig("config/development.yaml")
	if err != nil {
		log.Fatal(err)
	}

	app := &application{
		config: cfg,
	}

	mux := chi.NewRouter()
	app.run(mux)
}
