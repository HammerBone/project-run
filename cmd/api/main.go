package main

import (
	"log"

	"github.com/HammerBone/project-run/internal/config"
	"github.com/HammerBone/project-run/internal/database"
)

func main() {
	cfg, err := config.LoadConfig("config/development.yaml")
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewDBConnection(cfg.DB.DBAddr)
	if err != nil {
		log.Fatal(err)
	}

	server := NewServer(cfg, db)

	server.run()
}
