// internal/repository/db.go
package repository

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// SetupDatabase initializes the database connection and schema
func SetupDatabase(dbPath string) (*sql.DB, error) {
	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Check connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

// createTables sets up the database schema
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

	// Add other tables as needed

	return nil
}
