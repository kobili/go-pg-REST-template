package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	var direction = flag.String("dir", "up", "The direction to run the migrations: 'up' or 'down'")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	pgHost := os.Getenv("POSTGRES_HOST")
	pgPort := os.Getenv("POSTGRES_PORT")
	pgUser := os.Getenv("POSTGRES_USER")
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	pgDatabase := os.Getenv("POSTGRES_DATABASE")

	pgConnectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		pgUser,
		pgPassword,
		pgHost,
		pgPort,
		pgDatabase,
	)

	log.Printf("Migrating %swards on %s", *direction, pgConnectionString)

	m, err := migrate.New(
		"file://sql",
		pgConnectionString,
	)
	if err != nil {
		log.Fatal(err)
	}

	if *direction == "up" {
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := m.Steps(-1); err != nil {
			log.Fatal(err)
		}
	}
}
