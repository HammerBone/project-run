package main

import (
	"log"

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

	userStore := user.NewPostgreUserStorage(db)
	userService := user.NewUserService(cfg.App.JWTSecret, userStore)
	userHandler := user.NewUserHandler(userService)

	postStore := post.NewPostgrePostStorage(db)
	postService := post.NewPostService(postStore)
	postHandler := post.NewPostHandler(postService)

	server := NewServer(cfg, db, *userHandler, *postHandler)

	server.run()
}
