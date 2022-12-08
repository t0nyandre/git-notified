package postgres

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgres() (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres",
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_DB"),
			os.Getenv("POSTGRES_SSLMODE")))
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to database")

	return db, nil
}
