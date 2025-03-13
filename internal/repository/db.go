package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Postgres driver
	"github.com/mark-c-hall/pds-go/internal/config"
)

func SetupDatabase(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.Database.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	return db, nil
}

func createTables(db *sql.DB) error {
	// Accounts table
	accountsSchema := `
	CREATE TABLE IF NOT EXISTS accounts (
		did TEXT PRIMARY KEY,
		handle TEXT UNIQUE NOT NULL,
		email TEXT NOT NULL,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_accounts_handle ON accounts(handle);
	`

	_, err := db.Exec(accountsSchema)
	if err != nil {
		return err
	}

	return nil
}
