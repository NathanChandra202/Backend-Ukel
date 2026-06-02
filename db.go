package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func koneksiDatabase() error {
	// Muat backend/.env (abaikan jika tidak ada)
	_ = godotenv.Load()

	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "kontribid")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("sql.Open: %w", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	if err := DB.Ping(); err != nil {
		return fmt.Errorf(
			"tidak bisa menghubungi PostgreSQL (%s:%s, db=%s, user=%s): %w",
			host, port, dbname, user, err,
		)
	}

	log.Printf("Database OK → %s@%s:%s/%s", user, host, port, dbname)
	return nil
}
