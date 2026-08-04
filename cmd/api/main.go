package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/HammerBone/project-run/internal/config"
	"github.com/HammerBone/project-run/internal/database"
	"github.com/HammerBone/project-run/internal/post"
	"github.com/HammerBone/project-run/internal/user"
)

func main() {
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewDBConnection(cfg.DB.DBAddr)
	if err != nil {
		log.Fatal(err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	userStore := user.NewPostgreUserStorage(db)
	userService := user.NewUserService(cfg.App.JWTSecret, userStore)
	userHandler := user.NewUserHandler(logger, userService)

	postStore := post.NewPostgrePostStorage(db)
	postService := post.NewPostService(postStore)
	postHandler := post.NewPostHandler(logger, postService)

	server := NewServer(cfg, db, *userHandler, *postHandler)

	server.run()
}
