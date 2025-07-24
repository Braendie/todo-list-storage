package main

import (
	"flag"
	"log"

	"github.com/Braendie/todo-list-storage/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var migrationsTable string
	flag.StringVar(&migrationsTable, "migrations-table", "", "name of migrations")

	cfg := config.MustLoad()

	var migrationsPath string
	if migrationsTable == "migrations_test" {
		migrationsPath = cfg.MigrationsTestPath
	} else {
		migrationsPath = cfg.MigrationsPath
	}

	m, err := migrate.New("file://"+migrationsPath, cfg.StoragePostgresCon)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No new migrations to apply")
		} else {
			log.Fatalf("Failed to apply migrations: %v", err)
		}
	}
	
	log.Println("Migrations applied successfully")
}
