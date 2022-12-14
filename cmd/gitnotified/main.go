package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"github.com/t0nyandre/git-notified/internal/pkg/database/postgres"
	"github.com/t0nyandre/git-notified/pkg/logger"
)

func init() {
	if err := godotenv.Load("config/env/.env"); err != nil {
		log.Fatalf("Could not load .env file: %v", err)
	}
}

func main() {
	ctx := context.Background()
	router := chi.NewRouter()
	logger := logger.NewLogger()
	postgres, err := postgres.NewPostgres(logger)
	if err != nil {
		logger.Fatalw("Could not connect to database",
			"error", err)
	}

	// Add logger, postgres and router to the context
	ctx = context.WithValue(ctx, "logger", logger)
	ctx = context.WithValue(ctx, "router", router)
	ctx = context.WithValue(ctx, "postgres", postgres)

	// user.NewHandler(ctx, user.NewRepository(postgres))
	//    r.Mount("/user", oAuthGithub)
	// r.Get("/auth/github/login", oAuthGithub.GithubLogin)
	// r.Get("/auth/github/callback", oAuthGithub.GithubCallback)

	logger.Infow("Server successfully up and running", "host", os.Getenv("APP_HOST"), "port", os.Getenv("APP_PORT"))
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT")), router); err != nil {
		logger.Fatalw("Could not start server", "error", err)
	}
}
