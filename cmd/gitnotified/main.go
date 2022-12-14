package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"github.com/t0nyandre/git-notified/internal/auth/github"
	"github.com/t0nyandre/git-notified/internal/pkg/database/postgres"
	"github.com/t0nyandre/git-notified/internal/user"
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
	ctx = context.WithValue(ctx, "postgres", postgres)

	user := user.NewHandler(ctx)
	github := github.New(ctx, os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_CLIENT_SECRET"), os.Getenv("GITHUB_CALLBACK_URL"), "user:email")

	router.Mount("/auth", github.Routes(router))
	logger.Infow("Successfully added routes", "routes", "/auth")
	router.Mount("/user", user.Routes(router))
	logger.Infow("Successfully added routes", "routes", "/user")

	logger.Infow("Server successfully up and running", "host", os.Getenv("APP_HOST"), "port", os.Getenv("APP_PORT"))
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT")), router); err != nil {
		logger.Fatalw("Could not start server", "error", err)
	}
}
