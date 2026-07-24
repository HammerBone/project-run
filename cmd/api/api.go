package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/HammerBone/project-run/internal/config"
	"github.com/HammerBone/project-run/internal/user"
	"github.com/go-chi/chi/v5"
)

type server struct {
	router *chi.Mux
	config *config.Config
	db     *sql.DB

	userHandler user.UserHandler
}

func NewServer(cfg *config.Config, db *sql.DB, userHandler user.UserHandler) *server {
	return &server{
		router: chi.NewRouter(),
		config: cfg,
		db:     db,

		userHandler: userHandler,
	}
}

func (s *server) run() error {
	srv := &http.Server{
		Addr:    s.config.App.ServerPort,
		Handler: s.router,
	}

	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Project-run backend server is running"))
	})

	s.router.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/create", s.userHandler.CreateUser)
			r.Post("/login", s.userHandler.LoginUser)
			r.Post("/logout", s.userHandler.LogoutUser)
		})
	})

	log.Println("Listening to port", s.config.App.ServerPort)
	return srv.ListenAndServe()
}
