package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/HammerBone/project-run/internal/config"
	"github.com/go-chi/chi/v5"
)

type server struct {
	router chi.Router
	config *config.Config
	db *sql.DB
}

func NewServer(cfg *config.Config, db *sql.DB) *server {
	return &server{
		router: chi.NewRouter(),
		config: cfg,
		db: db,
	}
}

func (s *server) run() error {
	srv := &http.Server{
		Addr: s.config.App.ServerPort,
		Handler: s.router,
	}

	s.router.Get("/", func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Project-run backend server is running"))
	})
	
	log.Println("Listening to port", s.config.App.ServerPort)
	return srv.ListenAndServe()
}
