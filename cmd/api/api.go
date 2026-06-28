package main

import (
	"log"
	"net/http"

	"github.com/HammerBone/project-run/internal/config"
)

type application struct {
	config *config.Config
	dbConfig config.DBConfig
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr: app.config.App.ServerPort,
		Handler: mux,
	}
	
	log.Println("Listening to port", app.config.App.ServerPort)
	return srv.ListenAndServe()
}
