package db

import (
	"fmt"
	"goevent/internal/config"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func Connect(cfg *config.Config) (*sqlx.DB, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Connect: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	DB = db
	log.Println("Connected to Postgres")

	if err := applyMigrations(dbURL); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return db, nil
}

func applyMigrations(dbURL string) error {
	m, err := migrate.New("file://migration", dbURL)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	log.Println("Migrations applied (if any)")
	return nil
}
