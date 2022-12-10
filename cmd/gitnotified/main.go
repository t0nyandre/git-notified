package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"github.com/t0nyandre/git-notified/internal/auth/github"
	"github.com/t0nyandre/git-notified/internal/pkg/database/postgres"
	"github.com/t0nyandre/git-notified/pkg/logger"
)

func init() {
	if err := godotenv.Load("config/env/.env"); err != nil {
		log.Fatalf("Could not load .env file: %v", err)
	}
}

func main() {
	r := chi.NewRouter()
	_ = logger.NewLogger()

	oAuthGithub := github.NewGithub()
	// Connect to database
	_, err := postgres.NewPostgres()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	r.Get("/auth/github/login", oAuthGithub.GithubLogin)
	r.Get("/auth/github/callback", oAuthGithub.GithubCallback)

	log.Println(fmt.Sprintf("Server running on %s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT")))
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT")), r); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
