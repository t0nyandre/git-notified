package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("config/env/.env"); err != nil {
		log.Fatalf("Could not load .env file: %v", err)
	}
}

func main() {
	router := chi.NewRouter()

	fmt.Println(os.Getenv("THISISIT"))

	http.ListenAndServe(":4000", router)
}
