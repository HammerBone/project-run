package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/HammerBone/project-run/internal/config"
	"github.com/HammerBone/project-run/internal/middleware"
	"github.com/HammerBone/project-run/internal/post"
	"github.com/HammerBone/project-run/internal/user"
	"github.com/HammerBone/project-run/internal/util"
	"github.com/go-chi/chi/v5"
)

type server struct {
	router *chi.Mux
	config *config.Config
	db     *sql.DB

	userHandler user.UserHandler
	postHandler post.PostHandler
}

func NewServer(cfg *config.Config, db *sql.DB, userHandler user.UserHandler, postHandler post.PostHandler) *server {
	return &server{
		router: chi.NewRouter(),
		config: cfg,
		db:     db,

		userHandler: userHandler,
		postHandler: postHandler,
	}
}

func (s *server) run() error {
	srv := &http.Server{
		Addr:    s.config.App.ServerPort,
		Handler: s.router,
	}
	jwtGenerator := util.NewJWTGenerator(s.config.App.JWTSecret)

	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Project-run backend server is running"))
	})

	s.router.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", s.userHandler.CreateUser)
			r.Post("/login", s.userHandler.LoginUser)
			r.Post("/logout", s.userHandler.LogoutUser)
		})

		r.Route("/posts", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(jwtGenerator))
			r.Post("/", s.postHandler.CreatePost)
			r.Patch("/", s.postHandler.EditPost)
		})
	})

	log.Println("Listening to port", s.config.App.ServerPort)
	return srv.ListenAndServe()
}
