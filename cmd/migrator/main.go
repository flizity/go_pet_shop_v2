package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	var migrationsPath, migrationsTable string
	dbURL := os.Getenv("DATABASE_URL")

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "schema_migrations", "name of migrations table")
	flag.Parse()

	if dbURL == "" {
		log.Fatal("DATABASE_URL not set in environment")
	}
	if migrationsPath == "" {
		log.Fatal("migrations-path is required")
	}

	var separator string
	if strings.Contains(dbURL, "?") {
		separator = "&"
	} else {
		separator = "?"
	}
	databaseURL := fmt.Sprintf("%s%sx-migrations-table=%s", dbURL, separator, migrationsTable)
	sourceURL := "file://" + migrationsPath

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		log.Fatal(err)
	}

	fmt.Println("migrations applied")
}
