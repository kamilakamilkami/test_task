package config

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations() {
	dbURL := os.Getenv("DATABASE_URL")

	m, err := migrate.New("file://../migrations", dbURL)

	if err != nil {
		log.Fatalf("failed to initialize migrations: %v", err)
	}

	err = m.Up()
	if err != nil && err.Error() != "no change" {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	log.Println("✅ Migrations applied successfully")
}
