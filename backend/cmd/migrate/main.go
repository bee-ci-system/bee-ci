package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const defaultMigrationsPath = "migrations"

func main() {
	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	if command == "help" || command == "-h" || command == "--help" {
		printUsage()
		return
	}

	if command != "up" && command != "down" && command != "version" {
		printUsage()
		log.Fatalf("unknown command: %s", command)
	}

	db, err := sql.Open("postgres", postgresConnectionString())
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("error connecting to Postgres database: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	migrationsURL := "file://" + migrationsPath()
	m, err := migrate.NewWithDatabaseInstance(migrationsURL, "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		sourceErr, databaseErr := m.Close()
		if sourceErr != nil {
			log.Printf("error closing migration source: %v", sourceErr)
		}
		if databaseErr != nil {
			log.Printf("error closing migration database: %v", databaseErr)
		}
	}()

	switch command {
	case "up":
		runUp(m)
	case "down":
		runDown(m, os.Args[2:])
	case "version":
		printVersion(m)
	}
}

func runUp(m *migrate.Migrate) {
	err := m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("database already up to date")
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Println("database migrations applied")
}

func runDown(m *migrate.Migrate, args []string) {
	steps := 1
	if len(args) > 0 {
		parsedSteps, err := strconv.Atoi(args[0])
		if err != nil || parsedSteps < 1 {
			log.Fatalf("down step count must be a positive integer, got: %q", args[0])
		}
		steps = parsedSteps
	}

	err := m.Steps(-steps)
	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("database already at the first migration")
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("rolled back %d migration(s)", steps)
}

func printVersion(m *migrate.Migrate) {
	version, dirty, err := m.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		log.Println("database has no applied migrations")
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("database migration version: %d, dirty: %t", version, dirty)
}

func postgresConnectionString() string {
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		return databaseURL
	}

	dbHost := mustGetenv("DB_HOST")
	dbPort := mustGetenv("DB_PORT")
	dbUser := mustGetenv("DB_USER")
	dbPassword := mustGetenv("DB_PASSWORD")
	dbName := mustGetenv("DB_NAME")
	dbOpts := mustGetenv("DB_OPTS")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s %s", dbHost, dbPort, dbUser, dbPassword, dbName, dbOpts)
}

func migrationsPath() string {
	if path := os.Getenv("MIGRATIONS_PATH"); path != "" {
		return path
	}
	return defaultMigrationsPath
}

func mustGetenv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalf("%s env var is empty or not set", name)
	}
	return value
}

func printUsage() {
	log.Printf(`usage:
  migrate up
  migrate down [steps]
  migrate version

environment:
  DATABASE_URL      optional Postgres URL, used if set
  DB_HOST          required if DATABASE_URL is not set
  DB_PORT          required if DATABASE_URL is not set
  DB_USER          required if DATABASE_URL is not set
  DB_PASSWORD      required if DATABASE_URL is not set
  DB_NAME          required if DATABASE_URL is not set
  DB_OPTS          required if DATABASE_URL is not set, for example "sslmode=require"
  MIGRATIONS_PATH  optional path to migration files, default "migrations"`)
}
