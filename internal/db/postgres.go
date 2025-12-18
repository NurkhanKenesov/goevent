package db

import (
	"fmt"
	"goevent/internal/config"
	"log"
	"os"
	"path/filepath"
	"strings"

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
	// Determine migration source directory. The working directory when running
	// the app may be `cmd/app` or project root; check common relative paths
	// and use the first existing one.
	cwd, _ := os.Getwd()
	candidates := []string{"migration", "../migration", "../../migration"}
	var src string
	for _, c := range candidates {
		p := filepath.Join(cwd, c)
		if fi, err := os.Stat(p); err == nil && fi.IsDir() {
			abs, _ := filepath.Abs(p)
			absSlash := filepath.ToSlash(abs)
			// On Windows absolute paths start with a drive letter (C:/...)
			// file URL should be file:///C:/path whereas on Unix it's file:///path
			if os.PathSeparator == '\\' {
				src = "file:///" + absSlash
			} else {
				src = "file://" + absSlash
			}
			break
		}
	}
	if src == "" {
		// fallback to relative path (may still fail)
		src = "file://migration"
	}

	log.Println("Using migration source:", src)

	// If we found an absolute migration path, switch into it and use a relative file://. source
	var m *migrate.Migrate
	var err error
	var migAbs string
	// derive migAbs from src when possible
	if src != "file://migration" {
		if os.PathSeparator == '\\' {
			migAbs = strings.TrimPrefix(src, "file:///")
		} else {
			migAbs = strings.TrimPrefix(src, "file://")
		}
	}
	if migAbs != "" {
		prevWd, _ := os.Getwd()
		if err := os.Chdir(migAbs); err == nil {
			m, err = migrate.New("file://.", dbURL)
			_ = os.Chdir(prevWd)
			if err != nil {
				return err
			}
		} else {
			// fallback
			m, err = migrate.New(src, dbURL)
			if err != nil {
				return err
			}
		}
	} else {
		m, err = migrate.New(src, dbURL)
		if err != nil {
			return err
		}
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	log.Println("Migrations applied (if any)")
	return nil
}
